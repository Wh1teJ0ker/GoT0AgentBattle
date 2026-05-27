package battle

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

// EventName 是后端向前端广播房间事件时使用的统一事件名。
const EventName = "got0:broadcast"

// RoomConfigInput 描述前端创建房间时提交的配置参数。
type RoomConfigInput struct {
	Topic               string   `json:"topic"`
	Mode                string   `json:"mode"`
	AgentCount          int      `json:"agentCount"`
	Model               string   `json:"model"`
	Rounds              int      `json:"rounds"`
	ProviderPath        string   `json:"providerPath"`
	PreferRealLLM       bool     `json:"preferRealLLM"`
	PersonaIDs          []string `json:"personaIds"`
	JudgeModelProfileID string   `json:"judgeModelProfileId"`
}

// AudienceActionInput 描述观众插嘴、导播指令等外部干预输入。
type AudienceActionInput struct {
	Message       string `json:"message"`
	TargetAgentID string `json:"targetAgentId"`
	Kind          string `json:"kind"`
}

// ModeOption 描述一种可选辩论模式的展示信息。
type ModeOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// ModelOption 描述一个模型位的前端展示信息。
type ModelOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// ModelProfileSummary 是模型档案的精简展示结构，用于前端下拉或状态说明。
type ModelProfileSummary struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Model       string `json:"model"`
	Description string `json:"description"`
}

// PersonaConfig 描述一个人格模板的完整配置。
type PersonaConfig struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Role                  string `json:"role"`
	Avatar                string `json:"avatar"`
	Color                 string `json:"color"`
	Style                 string `json:"style"`
	Persona               string `json:"persona"`
	Tagline               string `json:"tagline"`
	Aggressive            int    `json:"aggressive"`
	Toxicity              int    `json:"toxicity"`
	Enabled               bool   `json:"enabled"`
	DefaultModelProfileID string `json:"defaultModelProfileId"`
	SystemPrompt          string `json:"systemPrompt"`
}

// AppSettings 表示本地持久化保存的应用设置。
type AppSettings struct {
	DefaultConfig        RoomConfigInput `json:"defaultConfig"`
	Personas             []PersonaConfig `json:"personas"`
	JudgeModelProfileID  string          `json:"judgeModelProfileId"`
	PreferredThemeFlavor string          `json:"preferredThemeFlavor"`
}

// BootstrapData 是前端首屏启动时需要的一次性初始化数据。
type BootstrapData struct {
	Modes              []ModeOption          `json:"modes"`
	Models             []ModelOption         `json:"models"`
	ModelProfiles      []ModelProfileSummary `json:"modelProfiles"`
	DefaultConfig      RoomConfigInput       `json:"defaultConfig"`
	Settings           AppSettings           `json:"settings"`
	SettingsPath       string                `json:"settingsPath"`
	State              RoomState             `json:"state"`
	RealProviderReady  bool                  `json:"realProviderReady"`
	RealProviderPath   string                `json:"realProviderPath"`
	RealProviderNotice string                `json:"realProviderNotice"`
}

// AgentState 表示某个角色在当前房间中的运行时状态。
type AgentState struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Role           string `json:"role"`
	Avatar         string `json:"avatar"`
	Color          string `json:"color"`
	Style          string `json:"style"`
	Persona        string `json:"persona"`
	Tagline        string `json:"tagline"`
	Aggressive     int    `json:"aggressive"`
	Toxicity       int    `json:"toxicity"`
	Enabled        bool   `json:"enabled"`
	Anger          int    `json:"anger"`
	Support        int    `json:"support"`
	TokenUsage     int    `json:"tokenUsage"`
	Momentum       int    `json:"momentum"`
	RoastCount     int    `json:"roastCount"`
	Status         string `json:"status"`
	Model          string `json:"model"`
	ModelProfileID string `json:"modelProfileId"`
	LastLine       string `json:"lastLine"`
	LastTarget     string `json:"lastTarget"`
	CurrentTurn    bool   `json:"currentTurn"`
}

// DebateMessage 表示聊天室里的一条消息。
type DebateMessage struct {
	ID             string   `json:"id"`
	AgentID        string   `json:"agentId"`
	AgentName      string   `json:"agentName"`
	AgentAvatar    string   `json:"agentAvatar"`
	Color          string   `json:"color"`
	Content        string   `json:"content"`
	ReplyTo        string   `json:"replyTo"`
	Mentions       []string `json:"mentions"`
	Timestamp      int64    `json:"timestamp"`
	Round          int      `json:"round"`
	Tone           string   `json:"tone"`
	IsJudge        bool     `json:"isJudge"`
	IsAudience     bool     `json:"isAudience"`
	HeatDelta      int      `json:"heatDelta"`
	SupportDelta   int      `json:"supportDelta"`
	Impact         int      `json:"impact"`
	Cue            string   `json:"cue"`
	PerformanceTag string   `json:"performanceTag"`
}

// JudgeSummary 表示裁判对本轮或整场辩论的总结结果。
type JudgeSummary struct {
	Round               int      `json:"round"`
	WinnerAgentID       string   `json:"winnerAgentId"`
	WinnerName          string   `json:"winnerName"`
	WinnerReason        string   `json:"winnerReason"`
	WinnerSide          string   `json:"winnerSide"`
	EffectiveArguments  []string `json:"effectiveArguments"`
	InvalidTrashTalk    []string `json:"invalidTrashTalk"`
	ShowmanshipScore    int      `json:"showmanshipScore"`
	EntertainmentRating string   `json:"entertainmentRating"`
}

// RoomState 表示房间在任意时刻的完整状态快照。
type RoomState struct {
	RoomID              string          `json:"roomId"`
	Topic               string          `json:"topic"`
	Mode                string          `json:"mode"`
	Model               string          `json:"model"`
	Status              string          `json:"status"`
	CurrentRound        int             `json:"currentRound"`
	TotalRounds         int             `json:"totalRounds"`
	Heat                int             `json:"heat"`
	DramaLevel          int             `json:"dramaLevel"`
	AudienceMood        string          `json:"audienceMood"`
	SupportLeader       string          `json:"supportLeader"`
	CurrentSpeaker      string          `json:"currentSpeaker"`
	CurrentTarget       string          `json:"currentTarget"`
	LastNotice          string          `json:"lastNotice"`
	AudienceQueue       int             `json:"audienceQueue"`
	StartedAt           int64           `json:"startedAt"`
	FinishedAt          int64           `json:"finishedAt"`
	Engine              string          `json:"engine"`
	ProviderPath        string          `json:"providerPath"`
	ProviderReady       bool            `json:"providerReady"`
	ProviderNotice      string          `json:"providerNotice"`
	JudgeModelProfileID string          `json:"judgeModelProfileId"`
	PersonaIDs          []string        `json:"personaIds"`
	Agents              []AgentState    `json:"agents"`
	Messages            []DebateMessage `json:"messages"`
	JudgeSummary        *JudgeSummary   `json:"judgeSummary"`
}

// EventPayload 是后端通过 Wails 事件广播给前端的统一消息结构。
type EventPayload struct {
	Type      string         `json:"type"`
	State     RoomState      `json:"state"`
	Message   *DebateMessage `json:"message,omitempty"`
	Summary   *JudgeSummary  `json:"summary,omitempty"`
	Notice    string         `json:"notice,omitempty"`
	Timestamp int64          `json:"timestamp"`
}

type agentPreset struct {
	ID         string
	Name       string
	Role       string
	Avatar     string
	Color      string
	Style      string
	Persona    string
	Tagline    string
	Aggressive int
	Toxicity   int
}

var modeOptions = []ModeOption{
	{ID: "free-for-all", Label: "自由辩论", Description: "所有人抢麦，谁 loud 谁上分。"},
	{ID: "red-vs-blue", Label: "红蓝对抗", Description: "按立场分组，输出更有条理。"},
	{ID: "host-mode", Label: "主持人模式", Description: "节奏更稳，适合做节目版演示。"},
	{ID: "king-of-hill", Label: "擂台模式", Description: "赢家留场，败者掉麦。"},
	{ID: "chaos", Label: "混战模式", Description: "全员乱斗，节目效果拉满。"},
}

var modelOptions = []ModelOption{
	{ID: "openai", Label: "OpenAI 模型位", Description: "预留给真实模型接入。"},
	{ID: "claude", Label: "Claude 模型位", Description: "预留给真实模型接入。"},
	{ID: "gemini", Label: "Gemini 模型位", Description: "预留给真实模型接入。"},
	{ID: "deepseek", Label: "DeepSeek 模型位", Description: "预留给真实模型接入。"},
	{ID: "ollama", Label: "Ollama 模型位", Description: "预留给本地模型接入。"},
}

func defaultConfig() RoomConfigInput {
	return RoomConfigInput{
		Topic:         "",
		Mode:          "chaos",
		AgentCount:    5,
		Model:         "openai",
		Rounds:        3,
		ProviderPath:  "",
		PreferRealLLM: true,
	}
}

func idleState() RoomState {
	settings := defaultSettings()
	cfg := settings.DefaultConfig
	agents := buildAgents(settings.Personas, cfg.PersonaIDs, cfg.AgentCount, cfg.Model)
	now := time.Now().UnixMilli()
	return RoomState{
		RoomID:              "",
		Topic:               cfg.Topic,
		Mode:                cfg.Mode,
		Model:               cfg.Model,
		Status:              "idle",
		CurrentRound:        0,
		TotalRounds:         cfg.Rounds,
		Heat:                8,
		DramaLevel:          18,
		AudienceMood:        "观众还在试麦，等第一轮开喷。",
		SupportLeader:       agents[0].Name,
		CurrentSpeaker:      "",
		CurrentTarget:       "",
		LastNotice:          "先建房，再让他们互喷。",
		AudienceQueue:       0,
		StartedAt:           0,
		FinishedAt:          0,
		Engine:              "待配置",
		ProviderPath:        "",
		ProviderReady:       false,
		ProviderNotice:      "未检测到 YAML 模型配置，当前无法启动辩论。",
		JudgeModelProfileID: cfg.JudgeModelProfileID,
		PersonaIDs:          append([]string(nil), cfg.PersonaIDs...),
		Agents:              agents,
		Messages: []DebateMessage{
			{
				ID:             "system-welcome",
				AgentID:        "host",
				AgentName:      "导播台",
				AgentAvatar:    "TV",
				Color:          "#f1f2f6",
				Content:        "欢迎来到 GoT0AgentBattle。今天的核心任务不是求真，是求节目效果。模型能力由本地 YAML 配置驱动。",
				Timestamp:      now,
				Round:          0,
				Tone:           "system",
				IsJudge:        false,
				IsAudience:     false,
				HeatDelta:      0,
				SupportDelta:   0,
				Impact:         12,
				Cue:            "片头暖场",
				PerformanceTag: "预热中",
			},
		},
		JudgeSummary: nil,
	}
}

func buildAgents(personas []PersonaConfig, personaIDs []string, count int, fallbackModel string) []AgentState {
	selected := selectPersonas(personas, personaIDs, count)
	if len(selected) == 0 {
		selected = selectPersonas(defaultPersonas(), nil, count)
	}
	count = int(math.Max(1, math.Min(float64(count), float64(len(selected)))))
	result := make([]AgentState, 0, count)
	for i := 0; i < count; i++ {
		item := selected[i]
		modelProfileID := item.DefaultModelProfileID
		if strings.TrimSpace(modelProfileID) == "" {
			modelProfileID = fallbackModel
		}
		result = append(result, AgentState{
			ID:             item.ID,
			Name:           item.Name,
			Role:           item.Role,
			Avatar:         item.Avatar,
			Color:          item.Color,
			Style:          item.Style,
			Persona:        item.Persona,
			Tagline:        item.Tagline,
			Aggressive:     item.Aggressive,
			Toxicity:       item.Toxicity,
			Enabled:        item.Enabled,
			Anger:          maxInt(18, item.Aggressive-24),
			Support:        40 + (i * 3),
			TokenUsage:     0,
			Momentum:       36 + (i * 5),
			RoastCount:     0,
			Status:         "待命",
			Model:          fallbackModel,
			ModelProfileID: modelProfileID,
			LastLine:       "还没开麦。",
			LastTarget:     "",
			CurrentTurn:    false,
		})
	}
	return result
}

func selectPersonas(personas []PersonaConfig, personaIDs []string, count int) []PersonaConfig {
	if len(personas) == 0 {
		return nil
	}
	index := make(map[string]PersonaConfig, len(personas))
	enabled := make([]PersonaConfig, 0, len(personas))
	for _, persona := range personas {
		index[persona.ID] = persona
		if persona.Enabled {
			enabled = append(enabled, persona)
		}
	}
	if len(enabled) == 0 {
		enabled = personas
	}

	selected := make([]PersonaConfig, 0, count)
	seen := map[string]struct{}{}
	for _, id := range personaIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		persona, ok := index[id]
		if !ok {
			continue
		}
		selected = append(selected, persona)
		seen[id] = struct{}{}
	}
	for _, persona := range enabled {
		if len(selected) >= count {
			break
		}
		if _, ok := seen[persona.ID]; ok {
			continue
		}
		selected = append(selected, persona)
		seen[persona.ID] = struct{}{}
	}
	return selected
}

func normalizeConfig(input RoomConfigInput) RoomConfigInput {
	cfg := input
	if strings.TrimSpace(cfg.Mode) == "" {
		cfg.Mode = defaultConfig().Mode
	}
	if strings.TrimSpace(cfg.Model) == "" {
		cfg.Model = defaultConfig().Model
	}
	if cfg.AgentCount <= 0 {
		cfg.AgentCount = defaultConfig().AgentCount
	}
	if cfg.Rounds <= 0 {
		cfg.Rounds = defaultConfig().Rounds
	}
	if cfg.Rounds > 6 {
		cfg.Rounds = 6
	}
	if cfg.AgentCount > 12 {
		cfg.AgentCount = 12
	}
	cfg.PreferRealLLM = true
	return cfg
}

func newRoomState(cfg RoomConfigInput, settings AppSettings, providerReady bool, providerPath string, providerNotice string) RoomState {
	cfg = normalizeConfig(cfg)
	now := time.Now().UnixMilli()
	agents := buildAgents(settings.Personas, cfg.PersonaIDs, cfg.AgentCount, cfg.Model)
	engine := "待配置"
	if providerReady {
		engine = "real"
	}
	return RoomState{
		RoomID:              fmt.Sprintf("room-%d", now),
		Topic:               cfg.Topic,
		Mode:                cfg.Mode,
		Model:               cfg.Model,
		Status:              "ready",
		CurrentRound:        0,
		TotalRounds:         cfg.Rounds,
		Heat:                12,
		DramaLevel:          24,
		AudienceMood:        "导播刚把题目打上屏，观众等着看谁先开团。",
		SupportLeader:       agents[0].Name,
		CurrentSpeaker:      "",
		CurrentTarget:       "",
		LastNotice:          "房间已创建，等待导播按下开战键。",
		AudienceQueue:       0,
		StartedAt:           0,
		FinishedAt:          0,
		Engine:              engine,
		ProviderPath:        providerPath,
		ProviderReady:       providerReady,
		ProviderNotice:      providerNotice,
		JudgeModelProfileID: cfg.JudgeModelProfileID,
		PersonaIDs:          append([]string(nil), cfg.PersonaIDs...),
		Agents:              agents,
		Messages: []DebateMessage{
			{
				ID:             fmt.Sprintf("room-created-%d", now),
				AgentID:        "host",
				AgentName:      "导播台",
				AgentAvatar:    "TV",
				Color:          "#f1f2f6",
				Content:        fmt.Sprintf("话题《%s》已载入，模式 %s，选手 %d 人。当前引擎：%s。", cfg.Topic, cfg.Mode, cfg.AgentCount, engine),
				Timestamp:      now,
				Round:          0,
				Tone:           "system",
				HeatDelta:      2,
				SupportDelta:   0,
				Impact:         18,
				Cue:            "房间建立",
				PerformanceTag: "建房成功",
			},
		},
	}
}

func (s RoomState) clone() RoomState {
	copyState := s
	copyState.Agents = append([]AgentState(nil), s.Agents...)
	copyState.Messages = append([]DebateMessage(nil), s.Messages...)
	if s.JudgeSummary != nil {
		summaryCopy := *s.JudgeSummary
		summaryCopy.EffectiveArguments = append([]string(nil), s.JudgeSummary.EffectiveArguments...)
		summaryCopy.InvalidTrashTalk = append([]string(nil), s.JudgeSummary.InvalidTrashTalk...)
		copyState.JudgeSummary = &summaryCopy
	}
	return copyState
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func clampInt(v int, minValue int, maxValue int) int {
	if v < minValue {
		return minValue
	}
	if v > maxValue {
		return maxValue
	}
	return v
}

func chooseDifferentAgentIndex(rng *rand.Rand, current int, total int) int {
	if total <= 1 {
		return current
	}
	offset := rng.Intn(total - 1)
	if offset >= current {
		offset++
	}
	return offset
}
