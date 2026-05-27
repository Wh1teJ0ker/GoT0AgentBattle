package battle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// runReal 执行完整的真实模型辩论流程。
// 单独放在这个文件，是为了把长流程编排和普通服务接口拆开，方便阅读与评审。
func (s *Service) runReal(ctx context.Context, providers DebateProviders, settings AppSettings, state RoomState) {
	select {
	case <-ctx.Done():
		s.logger.Logf("battle", "run skipped because context was cancelled before start")
		return
	default:
	}

	s.logger.Logf("battle", "starting real debate room=%s topic=%q agents=%d rounds=%d", state.RoomID, state.Topic, len(state.Agents), state.TotalRounds)
	s.emitRoundStart(1, "真实模型辩论开始，正在把所有嘴替拉进主舞台。")
	result, err := RunRealDebateWithHooks(state, settings, providers, func(round int, _ int, speaker AgentState) *AudienceActionInput {
		s.emitRoundStart(round, fmt.Sprintf("第 %d 回合开始，%s 准备接麦。", round, speaker.Name))
		return s.dequeueAudienceActionForTurn(speaker)
	})
	if err != nil {
		s.failRun("真实模型辩论失败，已中止。", err)
		return
	}

	for index, turn := range result.Turns {
		select {
		case <-ctx.Done():
			s.logger.Logf("battle", "run interrupted room=%s at turn=%d", state.RoomID, index+1)
			return
		default:
		}

		// 先发“正在思考”事件，再发实际消息，
		// 这样前端可以展示更完整的实时辩论节奏。
		s.emitRealWarmup(index, turn)
		time.Sleep(300 * time.Millisecond)
		payload := s.applyRealTurn(index, turn)
		s.emit(payload)
		time.Sleep(280 * time.Millisecond)
	}

	s.mu.Lock()
	if len(result.Turns) > 0 {
		s.room.JudgeSummary = &result.Summary
	}
	s.room.Status = "finished"
	s.room.CurrentSpeaker = ""
	s.room.CurrentRound = s.room.TotalRounds
	s.room.LastNotice = "真实模型辩论已结束，裁判已宣判。"
	s.room.FinishedAt = time.Now().UnixMilli()
	s.room.SupportLeader = supportLeaderName(s.room.Agents)
	for i := range s.room.Agents {
		s.room.Agents[i].CurrentTurn = false
		if s.room.Agents[i].Status == "思考中" {
			s.room.Agents[i].Status = "等待复盘"
		}
	}
	s.cancel = nil
	finalState := s.room.clone()
	s.mu.Unlock()

	if path, err := SaveTranscriptMarkdown(finalState); err == nil {
		finalState.LastNotice = fmt.Sprintf("真实模型辩论已结束，记录已保存到 %s", path)
		s.mu.Lock()
		s.room.LastNotice = finalState.LastNotice
		s.mu.Unlock()
		s.logger.Logf("battle", "room=%s finished transcript=%s winner=%s", finalState.RoomID, path, summaryWinner(finalState.JudgeSummary))
	} else {
		finalState.LastNotice = fmt.Sprintf("真实模型辩论已结束，但 Markdown 记录保存失败：%v", err)
		s.mu.Lock()
		s.room.LastNotice = finalState.LastNotice
		s.mu.Unlock()
		s.logger.Logf("battle", "room=%s finished transcript save failed: %v", finalState.RoomID, err)
	}

	s.emit(EventPayload{
		Type:      "room.finished",
		State:     finalState,
		Summary:   finalState.JudgeSummary,
		Notice:    finalState.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
}

// emitRoundStart 更新房间回合信息，并广播新一轮开始事件。
func (s *Service) emitRoundStart(round int, notice string) {
	s.mu.Lock()
	s.room.CurrentRound = round
	s.room.LastNotice = notice
	s.room.AudienceMood = roomAudienceMood(s.room.Heat, s.room.DramaLevel, s.room.AudienceQueue)
	state := s.room.clone()
	s.mu.Unlock()

	s.emit(EventPayload{
		Type:      "round.started",
		State:     state,
		Notice:    notice,
		Timestamp: time.Now().UnixMilli(),
	})
}

// emitRealWarmup 在选手正式输出前，把当前火力点切到对应选手。
func (s *Service) emitRealWarmup(index int, turn RealDebateTurn) {
	s.mu.Lock()
	for i := range s.room.Agents {
		s.room.Agents[i].CurrentTurn = false
		if s.room.Agents[i].Status == "思考中" {
			s.room.Agents[i].Status = "等着反咬"
		}
	}
	s.room.CurrentSpeaker = turn.Stage.SpeakerLabel
	s.room.CurrentTarget = turn.Message.ReplyTo
	s.room.LastNotice = buildWarmupNotice(turn)
	if turn.Action != nil {
		s.room.AudienceMood = fmt.Sprintf("场外刚刷了一条“%s”，全场等着看 %s 怎么接。", turn.Action.Message, turn.Stage.SpeakerLabel)
	}
	for i := range s.room.Agents {
		if s.room.Agents[i].ID == turn.Stage.SpeakerID {
			s.room.Agents[i].CurrentTurn = true
			s.room.Agents[i].Status = "思考中"
		}
		if s.room.Agents[i].Name == turn.Message.ReplyTo {
			s.room.Agents[i].Status = "被点名"
		}
	}
	state := s.room.clone()
	s.mu.Unlock()

	s.emit(EventPayload{
		Type:      "speaker.thinking",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli() + int64(index),
	})
}

// applyRealTurn 是单次发言落地的唯一入口。
// 这里统一修改选手状态、消息流、热度和支持率，避免副作用分散在多个函数里。
func (s *Service) applyRealTurn(index int, turn RealDebateTurn) EventPayload {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := turn.Message
	msg.ID = fmt.Sprintf("real-%d-%d", index, time.Now().UnixNano())
	msg.Timestamp = time.Now().UnixMilli()
	msg.Round = extractRoundFromStage(turn.Stage.Title, minInt(index+1, s.room.TotalRounds))
	impact := estimateTurnImpact(turn)
	msg.Impact = impact
	msg.Cue = turnCue(turn, impact)
	msg.HeatDelta = impact / 7
	msg.SupportDelta = maxInt(2, impact/10)

	for idx := range s.room.Agents {
		agent := s.room.Agents[idx]
		if agent.ID != turn.Stage.SpeakerID {
			if agent.Name == turn.Message.ReplyTo {
				agent.Anger = clampInt(agent.Anger+maxInt(5, impact/9), 0, 100)
				agent.Momentum = clampInt(agent.Momentum-impact/12, 0, 100)
				agent.Status = "被当场点名"
			} else if agent.Status == "思考中" || agent.Status == "被点名" {
				agent.Status = "观察战局"
			}
			s.room.Agents[idx] = agent
			continue
		}

		agent.Name = turn.Stage.SpeakerLabel
		agent.Role = turn.Stage.Role
		agent.Model = turn.Stage.ModelName
		agent.ModelProfileID = turn.Stage.ModelProfileID
		agent.LastLine = turn.Stage.Output
		agent.LastTarget = turn.Message.ReplyTo
		agent.TokenUsage += len([]rune(turn.Stage.Output)) * 2
		agent.Support = clampInt(agent.Support+msg.SupportDelta, 0, 100)
		agent.Anger = clampInt(agent.Anger+maxInt(3, impact/11), 0, 100)
		agent.Momentum = clampInt(agent.Momentum+maxInt(4, impact/8), 0, 100)
		agent.RoastCount++
		agent.Status = cueStatus(msg.Cue, turn.Stage.Role)
		agent.CurrentTurn = false
		s.room.Agents[idx] = agent
		msg.AgentID = agent.ID
		msg.AgentName = agent.Name
		msg.AgentAvatar = agent.Avatar
		msg.Color = agent.Color
		break
	}

	s.room.Messages = append(s.room.Messages, msg)
	s.room.Heat = clampInt(s.room.Heat+maxInt(4, msg.HeatDelta), 0, 100)
	s.room.DramaLevel = clampInt(s.room.DramaLevel+maxInt(5, impact/8), 0, 100)
	s.room.AudienceMood = roomAudienceMood(s.room.Heat, s.room.DramaLevel, s.room.AudienceQueue)
	s.room.SupportLeader = supportLeaderName(s.room.Agents)
	s.room.CurrentSpeaker = turn.Stage.SpeakerLabel
	s.room.CurrentTarget = turn.Message.ReplyTo
	s.room.LastNotice = fmt.Sprintf("%s 已经开火，目标直指 %s。", turn.Stage.SpeakerLabel, fallbackText(turn.Message.ReplyTo, "全场"))
	if turn.Stage.Role == "judge" {
		summary, ok := inferSummaryFromJudgeMessage(turn.Message.Content, s.room.JudgeSummary)
		if ok {
			s.room.JudgeSummary = summary
		}
	}

	state := s.room.clone()
	return EventPayload{
		Type:      "message.created",
		State:     state,
		Message:   &msg,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	}
}

// inferSummaryFromJudgeMessage 当前直接复用已有裁判总结。
// 这个钩子保留下来，便于后续扩展成从裁判原文中二次提取摘要。
func inferSummaryFromJudgeMessage(_ string, current *JudgeSummary) (*JudgeSummary, bool) {
	if current == nil {
		return nil, false
	}
	return current, true
}

// failRun 在真实模型链路失败时统一收口状态、日志和前端事件。
func (s *Service) failRun(notice string, err error) {
	s.mu.Lock()
	s.room.Status = "stopped"
	s.room.CurrentSpeaker = ""
	s.room.LastNotice = fmt.Sprintf("%s %v", notice, err)
	s.cancel = nil
	state := s.room.clone()
	s.mu.Unlock()

	s.logger.Logf("battle", "room failed room=%s error=%v", state.RoomID, err)
	s.emit(EventPayload{
		Type:      "room.failed",
		State:     state,
		Notice:    state.LastNotice,
		Timestamp: time.Now().UnixMilli(),
	})
}

// emit 只负责把已经准备好的事件广播给前端。
// 所有状态修改都应当在调用前完成，这样事件发送本身保持无副作用。
func (s *Service) emit(payload EventPayload) {
	s.mu.Lock()
	ctx := s.ctx
	s.mu.Unlock()
	if ctx == nil {
		return
	}
	runtime.EventsEmit(ctx, EventName, payload)
}

// summaryWinner 提取裁判总结里的赢家名，主要给日志和提示文案使用。
func summaryWinner(summary *JudgeSummary) string {
	if summary == nil || summary.WinnerName == "" {
		return "unknown"
	}
	return summary.WinnerName
}

func estimateTurnImpact(turn RealDebateTurn) int {
	if turn.Stage.Role == "judge" {
		return 100
	}

	lengthScore := minInt(34, len([]rune(strings.TrimSpace(turn.Stage.Output)))/7)
	mentionScore := minInt(18, len(turn.Message.Mentions)*8)
	actionScore := 0
	if turn.Action != nil {
		actionScore = 16
	}
	replyScore := 0
	if strings.TrimSpace(turn.Message.ReplyTo) != "" {
		replyScore = 12
	}
	return clampInt(32+lengthScore+mentionScore+replyScore+actionScore, 20, 96)
}

func turnCue(turn RealDebateTurn, impact int) string {
	if turn.Stage.Role == "judge" {
		return "终局宣判"
	}
	if turn.Action != nil && strings.TrimSpace(turn.Action.TargetAgentID) != "" {
		return "精准点杀"
	}
	if impact >= 84 {
		return "全场沸腾"
	}
	if impact >= 70 {
		return "火力拉满"
	}
	if strings.Contains(turn.Stage.Output, "@") {
		return "公开点名"
	}
	return "持续施压"
}

func cueStatus(cue string, role string) string {
	if role == "judge" {
		return "宣判完成"
	}
	switch cue {
	case "精准点杀":
		return "点名爆破"
	case "全场沸腾":
		return "火力失控"
	case "火力拉满":
		return "压着输出"
	case "公开点名":
		return "当面拆台"
	default:
		return "已输出"
	}
}

func roomAudienceMood(heat int, drama int, queue int) string {
	switch {
	case queue >= 3:
		return "弹幕已经开始堆栈，观众明显在拱火。"
	case drama >= 78:
		return "观众情绪已经被点燃，整个场子都在等下一次反击。"
	case heat >= 66:
		return "节奏正在上扬，场外开始有人站队。"
	default:
		return "观众还在观察哪边会先翻车。"
	}
}

func buildWarmupNotice(turn RealDebateTurn) string {
	if turn.Action != nil {
		return fmt.Sprintf("%s 正在接住场外节奏，准备冲着 %s 开火。", turn.Stage.SpeakerLabel, fallbackText(turn.Message.ReplyTo, "全场"))
	}
	return fmt.Sprintf("%s 正在锁定 %s，准备生成下一段火力输出。", turn.Stage.SpeakerLabel, fallbackText(turn.Message.ReplyTo, "全场"))
}

func extractRoundFromStage(title string, fallback int) int {
	var round int
	if _, err := fmt.Sscanf(title, "Round %d -", &round); err == nil && round > 0 {
		return round
	}
	return fallback
}
