package battle

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Service 是桌面端辩论服务的核心协调器。
// 它负责持有房间运行时状态、设置、模型配置以及对前端的事件广播能力。
type Service struct {
	mu                   sync.Mutex
	ctx                  context.Context
	rng                  *rand.Rand
	logger               *Logger
	room                 RoomState
	settings             AppSettings
	settingsPath         string
	pendingActions       []AudienceActionInput
	cancel               context.CancelFunc
	providers            DebateProviders
	providerPath         string
	providerReady        bool
	providerNotice       string
	realEngineConfigured bool
}

// NewService 创建核心业务服务。
// 这个服务会长期驻留在桌面应用进程中，因此同时维护设置、房间状态、模型配置和本地日志。
func NewService() *Service {
	settings, loadedSettingsPath, err := LoadSettings()
	if err != nil {
		settings = defaultSettings()
		loadedSettingsPath = settingsPath()
	}
	s := &Service{
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:       NewLogger(),
		settings:     settings,
		settingsPath: loadedSettingsPath,
		room:         idleState(),
	}
	s.room = newRoomState(s.settings.DefaultConfig, s.settings, false, "", "未检测到 YAML 模型配置，当前无法启动辩论。")
	s.room.Status = "idle"
	s.room.RoomID = ""
	s.room.LastNotice = "先准备 YAML 模型配置，再建房开战。"
	s.room.Messages = idleState().Messages
	s.reloadProviders()
	s.logger.Logf("service", "service initialized settings=%s log=%s", s.settingsPath, s.logger.Path())
	return s
}

// reloadProviders 重新加载本地模型配置文件。
// 每次建房前都会调用它，以确保当前运行使用的是最新的 YAML 配置。
func (s *Service) reloadProviders() {
	cfg, path, err := LoadProviderConfig()
	if err != nil {
		s.providers = DebateProviders{}
		s.providerPath = ""
		s.providerReady = false
		s.realEngineConfigured = false
		s.providerNotice = "未检测到 YAML 模型配置，当前无法启动辩论。"
		s.room.ProviderReady = false
		s.room.ProviderPath = ""
		s.room.ProviderNotice = s.providerNotice
		s.room.Engine = "待配置"
		s.logger.Logf("config", "provider config unavailable: %v", err)
		return
	}
	s.providers = cfg
	s.providerPath = path
	s.providerReady = true
	s.realEngineConfigured = true
	s.providerNotice = "已检测到 YAML 模型配置，可以启动真实模型辩论。"
	s.room.ProviderReady = true
	s.room.ProviderPath = path
	s.room.ProviderNotice = s.providerNotice
	s.logger.Logf("config", "provider config loaded path=%s profiles=%d", path, len(cfg.Profiles))
}

// BindContext 绑定 Wails 运行时上下文。
// 后续事件广播依赖这个上下文把房间更新推送给前端。
func (s *Service) BindContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
	s.logger.Logf("service", "wails context bound")
}

// Bootstrap 返回前端首屏渲染需要的完整初始化数据。
// 这里一次性返回，避免前端启动后再发起多次拆散请求。
func (s *Service) Bootstrap() BootstrapData {
	s.mu.Lock()
	defer s.mu.Unlock()
	return BootstrapData{
		Modes:              append([]ModeOption(nil), modeOptions...),
		Models:             append([]ModelOption(nil), modelOptions...),
		ModelProfiles:      s.providers.Summaries(),
		DefaultConfig:      s.settings.DefaultConfig,
		Settings:           cloneSettings(s.settings),
		SettingsPath:       s.settingsPath,
		State:              s.room.clone(),
		RealProviderReady:  s.providerReady,
		RealProviderPath:   s.providerPath,
		RealProviderNotice: s.providerNotice,
	}
}

// State 返回当前房间状态快照。
func (s *Service) State() RoomState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.room.clone()
}

// CreateRoom 根据前端输入创建新的房间状态。
// 它会重置旧运行、刷新模型配置，并把缺省值补齐到当前房间。
func (s *Service) CreateRoom(input RoomConfigInput) RoomState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		// 创建新房间后，旧房间遗留的运行协程就不再合法，必须先中断。
		s.cancel()
		s.cancel = nil
	}
	s.reloadProviders()
	s.pendingActions = nil
	input.ProviderPath = s.providerPath
	if len(input.PersonaIDs) == 0 {
		input.PersonaIDs = append([]string(nil), s.settings.DefaultConfig.PersonaIDs...)
	}
	if strings.TrimSpace(input.JudgeModelProfileID) == "" {
		input.JudgeModelProfileID = s.settings.JudgeModelProfileID
	}
	s.room = newRoomState(input, s.settings, s.providerReady, s.providerPath, s.providerNotice)
	state := s.room.clone()
	s.logger.Logf("room", "created room=%s topic=%q mode=%s personas=%d", state.RoomID, state.Topic, state.Mode, len(state.Agents))
	go s.emit(EventPayload{
		Type:      "room.created",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
	return state
}

// Start 启动当前房间的真实模型辩论流程。
// 只有在房间已创建且本地模型配置可用时才允许真正开战。
func (s *Service) Start() (RoomState, error) {
	s.mu.Lock()
	if s.room.RoomID == "" {
		s.mu.Unlock()
		return RoomState{}, fmt.Errorf("room has not been created")
	}
	if s.room.Status == "running" {
		state := s.room.clone()
		s.mu.Unlock()
		return state, nil
	}
	if !s.shouldUseRealEngineLocked() {
		s.mu.Unlock()
		return RoomState{}, fmt.Errorf("未检测到可用的 YAML 模型配置，无法启动辩论")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.room.Status = "running"
	s.room.StartedAt = time.Now().UnixMilli()
	s.room.Engine = "real"
	s.room.LastNotice = "真实模型辩论已启动，开始逐位出牌。"
	state := s.room.clone()
	providers := s.providers
	settings := cloneSettings(s.settings)
	s.mu.Unlock()

	go s.emit(EventPayload{
		Type:      "room.started",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})

	go s.runReal(ctx, providers, settings, state)
	return state, nil
}

// Stop 手动中断当前辩论，并把选手状态恢复到已暂停态。
func (s *Service) Stop() RoomState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.room.Status = "stopped"
	s.room.CurrentSpeaker = ""
	s.room.LastNotice = "已手动中断辩论。"
	for i := range s.room.Agents {
		s.room.Agents[i].CurrentTurn = false
		if s.room.Agents[i].Status != "待命" {
			s.room.Agents[i].Status = "被导播消音"
		}
	}
	state := s.room.clone()
	s.logger.Logf("room", "stopped room=%s", state.RoomID)
	go s.emit(EventPayload{
		Type:      "room.stopped",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
	return state
}

// SendAudienceAction 是给前端保留的语义化入口，内部直接复用队列逻辑。
func (s *Service) SendAudienceAction(input AudienceActionInput) (RoomState, error) {
	return s.QueueAudienceAction(input)
}

// QueueAudienceAction 把观众消息压入待处理队列，并立即回显到聊天流中。
func (s *Service) QueueAudienceAction(input AudienceActionInput) (RoomState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.room.RoomID == "" {
		return RoomState{}, fmt.Errorf("room has not been created")
	}
	message := strings.TrimSpace(input.Message)
	if message == "" {
		return RoomState{}, fmt.Errorf("message cannot be empty")
	}
	action := AudienceActionInput{
		Message:       message,
		TargetAgentID: strings.TrimSpace(input.TargetAgentID),
		Kind:          strings.TrimSpace(input.Kind),
	}
	if action.Kind == "" {
		action.Kind = "弹幕"
	}
	s.pendingActions = append(s.pendingActions, action)
	s.room.AudienceQueue = len(s.pendingActions)
	// 这里先把观众动作回显到前端，保证聊天区即时有反馈，
	// 真正影响模型输出则留给下一次回合调度消费。
	feedMessage := DebateMessage{
		ID:             fmt.Sprintf("aud-%d", time.Now().UnixNano()),
		AgentID:        "audience",
		AgentName:      "场外观众",
		AgentAvatar:    "DM",
		Color:          "#dfe4ea",
		Content:        fmt.Sprintf("[%s] %s", action.Kind, action.Message),
		Timestamp:      time.Now().UnixMilli(),
		Round:          s.room.CurrentRound,
		Tone:           "audience",
		IsAudience:     true,
		HeatDelta:      4,
		SupportDelta:   0,
		PerformanceTag: "插嘴成功",
	}
	s.room.Messages = append(s.room.Messages, feedMessage)
	s.room.Heat = clampInt(s.room.Heat+4, 0, 100)
	s.room.LastNotice = "场外弹幕已插入，下一位选手会被带节奏。"
	state := s.room.clone()
	s.logger.Logf("audience", "queued action room=%s kind=%s target=%s", state.RoomID, action.Kind, action.TargetAgentID)
	go s.emit(EventPayload{
		Type:      "audience.injected",
		State:     state,
		Message:   &feedMessage,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
	return state, nil
}

// shouldUseRealEngineLocked 判断当前是否具备真实模型启动条件。
// 调用方必须已经持有互斥锁。
func (s *Service) shouldUseRealEngineLocked() bool {
	return s.providerReady && s.realEngineConfigured
}

// Settings 返回当前内存中的设置快照。
func (s *Service) Settings() AppSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return cloneSettings(s.settings)
}

// SaveSettings 持久化设置，并把房间状态重置为新的默认配置。
func (s *Service) SaveSettings(input AppSettings) (AppSettings, error) {
	normalized, path, err := SaveSettings(input)
	if err != nil {
		return AppSettings{}, err
	}

	s.mu.Lock()
	s.settings = normalized
	s.settingsPath = path
	s.room = newRoomState(s.settings.DefaultConfig, s.settings, s.providerReady, s.providerPath, s.providerNotice)
	s.room.Status = "idle"
	s.room.RoomID = ""
	s.room.LastNotice = "设置已保存，等待下一场开喷。"
	s.pendingActions = nil
	state := s.room.clone()
	s.mu.Unlock()
	s.logger.Logf("settings", "saved settings path=%s personas=%d", path, len(normalized.Personas))

	go s.emit(EventPayload{
		Type:      "settings.saved",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
	return cloneSettings(normalized), nil
}

// GenerateRandomPersona 基于当前人格库生成一个新的随机人格模板。
func (s *Service) GenerateRandomPersona() PersonaConfig {
	s.mu.Lock()
	defer s.mu.Unlock()
	persona := GenerateRandomPersona(s.settings.Personas, s.settings.JudgeModelProfileID)
	s.logger.Logf("persona", "generated random persona id=%s name=%q", persona.ID, persona.Name)
	return persona
}

func supportLeaderName(agents []AgentState) string {
	if len(agents) == 0 {
		return ""
	}
	best := agents[0]
	for _, agent := range agents[1:] {
		if agent.Support > best.Support {
			best = agent
		}
	}
	return best.Name
}

// dequeueAudienceActionForTurn 在每次发言前取出一条待处理观众干预。
// 如果消息带了指定攻击对象，会优先把这条指令交给对应回合消化。
func (s *Service) dequeueAudienceActionForTurn(speaker AgentState) *AudienceActionInput {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pendingActions) == 0 {
		s.room.AudienceQueue = 0
		return nil
	}

	index := 0
	for i, action := range s.pendingActions {
		if strings.TrimSpace(action.TargetAgentID) == "" {
			continue
		}
		if action.TargetAgentID == speaker.ID {
			continue
		}
		index = i
		break
	}

	action := s.pendingActions[index]
	s.pendingActions = append(s.pendingActions[:index], s.pendingActions[index+1:]...)
	s.room.AudienceQueue = len(s.pendingActions)
	copyAction := action
	return &copyAction
}
