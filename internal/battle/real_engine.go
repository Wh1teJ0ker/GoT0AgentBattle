package battle

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// RealDebateTurn 表示真实模型引擎产出的单次发言及其阶段信息。
type RealDebateTurn struct {
	Stage   DebateStage
	Message DebateMessage
	Action  *AudienceActionInput
}

// RealDebateResult 表示一整场真实模型辩论执行后的完整结果。
type RealDebateResult struct {
	Turns    []RealDebateTurn
	Summary  JudgeSummary
	Judgment DebateJudgment
}

// RunRealDebate 使用真实模型驱动多人格辩论并生成最终裁判结果。
func RunRealDebate(state RoomState, settings AppSettings, providers DebateProviders) (RealDebateResult, error) {
	return RunRealDebateWithHooks(state, settings, providers, nil)
}

// RunRealDebateWithHooks 按回合执行真实模型辩论，并允许调度层在每次发言前注入外部干预。
func RunRealDebateWithHooks(
	state RoomState,
	settings AppSettings,
	providers DebateProviders,
	beforeTurn func(round int, speakerIndex int, speaker AgentState) *AudienceActionInput,
) (RealDebateResult, error) {
	memory := NewMemory(state.Topic)
	turns := make([]RealDebateTurn, 0, len(state.Agents)*maxInt(1, state.TotalRounds)+1)
	personaIndex := mapPersonas(settings.Personas)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for round := 1; round <= maxInt(1, state.TotalRounds); round++ {
		for speakerIndex, speaker := range state.Agents {
			var action *AudienceActionInput
			if beforeTurn != nil {
				action = beforeTurn(round, speakerIndex, speaker)
			}

			target := chooseRealTarget(rng, state.Agents, speakerIndex, action)
			persona, ok := personaIndex[speaker.ID]
			if !ok {
				persona = PersonaConfig{
					ID:                    speaker.ID,
					Name:                  speaker.Name,
					Role:                  speaker.Role,
					Style:                 speaker.Style,
					Persona:               speaker.Persona,
					Tagline:               speaker.Tagline,
					Aggressive:            speaker.Aggressive,
					Toxicity:              speaker.Toxicity,
					DefaultModelProfileID: speaker.ModelProfileID,
					SystemPrompt:          fmt.Sprintf("你是 %s，一个%s风格的%s。", speaker.Name, speaker.Persona, speaker.Role),
				}
			}

			profileID := speaker.ModelProfileID
			if strings.TrimSpace(profileID) == "" {
				profileID = persona.DefaultModelProfileID
			}
			provider, ok := providers.GetProfile(profileID)
			if !ok {
				return RealDebateResult{}, fmt.Errorf("missing model profile for %s: %s", speaker.Name, profileID)
			}

			prompt := buildPersonaPrompt(state.Topic, state.Mode, round, speaker, target, memory.Snapshot(), action, persona)
			output, err := ChatCompletion(provider, []chatMessage{
				{Role: "system", Content: prompt},
				{Role: "user", Content: fmt.Sprintf("现在轮到 %s 开麦，狠狠干碎对方。", speaker.Name)},
			})
			if err != nil {
				return RealDebateResult{}, fmt.Errorf("real turn failed for %s: %w", speaker.Name, err)
			}

			stage := DebateStage{
				Title:          fmt.Sprintf("Round %d - %s", round, speaker.Name),
				Role:           speaker.Role,
				SpeakerLabel:   speaker.Name,
				SpeakerID:      speaker.ID,
				ModelName:      provider.Model,
				ModelProfileID: profileID,
				Requirement:    "围绕主题、针对对手、兼顾节目效果地完成多人混战发言。",
				Output:         output,
			}
			memory.Add(stage)
			stage.MemoryAfter = memory.Snapshot()

			turns = append(turns, RealDebateTurn{
				Stage: stage,
				Message: DebateMessage{
					AgentID:        speaker.ID,
					AgentName:      speaker.Name,
					AgentAvatar:    speaker.Avatar,
					Color:          speaker.Color,
					Content:        output,
					ReplyTo:        target.Name,
					Mentions:       []string{target.Name},
					Round:          round,
					Tone:           speaker.Style,
					PerformanceTag: "真实模型输出",
				},
				Action: action,
			})
		}
	}

	judgeProfileID := state.JudgeModelProfileID
	if strings.TrimSpace(judgeProfileID) == "" {
		judgeProfileID = settings.JudgeModelProfileID
	}
	if strings.TrimSpace(judgeProfileID) == "" {
		judgeProfileID = providers.DefaultJudgeProfileID
	}
	judgeProvider, ok := providers.GetProfile(judgeProfileID)
	if !ok {
		return RealDebateResult{}, fmt.Errorf("missing judge model profile: %s", judgeProfileID)
	}

	judgeRaw, err := ChatCompletion(judgeProvider, []chatMessage{
		{Role: "system", Content: buildJudgePrompt(state.Topic, state.TotalRounds, memory.FullLog(), state.Agents)},
		{Role: "user", Content: "请直接给出这一整场的节目化裁决 JSON。"},
	})
	if err != nil {
		return RealDebateResult{}, fmt.Errorf("judge failed: %w", err)
	}

	judgment, err := parseJudgment(judgeRaw)
	if err != nil {
		return RealDebateResult{}, fmt.Errorf("failed to parse judge result: %w", err)
	}

	judgeStage := DebateStage{
		Title:          "Judge Final Adjudication",
		Role:           "judge",
		SpeakerLabel:   "Judge.exe",
		SpeakerID:      "judge",
		ModelName:      judgeProvider.Model,
		ModelProfileID: judgeProfileID,
		Requirement:    judgeRequirement(),
		Output:         judgeRaw,
	}
	memory.Add(judgeStage)
	judgeStage.MemoryAfter = memory.Snapshot()

	turns = append(turns, RealDebateTurn{
		Stage: judgeStage,
		Message: DebateMessage{
			AgentID:        "judge",
			AgentName:      "Judge.exe",
			AgentAvatar:    "JG",
			Color:          "#feca57",
			Content:        compactJudgeMessage(judgment),
			Tone:           "judge",
			IsJudge:        true,
			PerformanceTag: "最终裁决",
		},
	})

	return RealDebateResult{
		Turns:    turns,
		Summary:  judgmentToSummary(judgment),
		Judgment: judgment,
	}, nil
}

func mapPersonas(personas []PersonaConfig) map[string]PersonaConfig {
	index := make(map[string]PersonaConfig, len(personas))
	for _, persona := range personas {
		index[persona.ID] = persona
	}
	return index
}

func chooseRealTarget(rng *rand.Rand, agents []AgentState, speakerIndex int, action *AudienceActionInput) AgentState {
	if len(agents) == 0 {
		return AgentState{}
	}
	if action != nil && strings.TrimSpace(action.TargetAgentID) != "" {
		for _, agent := range agents {
			if agent.ID == action.TargetAgentID && agent.ID != agents[speakerIndex].ID {
				return agent
			}
		}
	}
	targetIndex := chooseDifferentAgentIndex(rng, speakerIndex, len(agents))
	return agents[targetIndex]
}

func judgmentToSummary(judgment DebateJudgment) JudgeSummary {
	winnerName := strings.TrimSpace(judgment.WinnerSide)
	winnerID := strings.ToLower(strings.ReplaceAll(winnerName, " ", "-"))
	if winnerID == "" {
		winnerID = "winner"
	}

	effective := append([]string(nil), judgment.KeyClashes...)
	if len(effective) > 3 {
		effective = effective[:3]
	}

	invalid := make([]string, 0, 3)
	scores := append([]SpeakerScore(nil), judgment.SpeakerScores...)
	sort.Slice(scores, func(i int, j int) bool {
		return scores[i].Breakdown.Total < scores[j].Breakdown.Total
	})
	for i := 0; i < len(scores) && i < 2; i++ {
		invalid = append(invalid, scores[i].Speaker+"："+scores[i].Comment)
	}

	showmanship := 72
	if value, ok := judgment.SideTotals["showmanship"]; ok && value > 0 {
		showmanship = clampInt(value, 0, 100)
	} else if len(judgment.KeyClashes) >= 3 {
		showmanship = 84
	}

	return JudgeSummary{
		Round:               1,
		WinnerAgentID:       winnerID,
		WinnerName:          winnerName,
		WinnerReason:        judgment.WinnerReason,
		WinnerSide:          winnerName,
		EffectiveArguments:  effective,
		InvalidTrashTalk:    invalid,
		ShowmanshipScore:    showmanship,
		EntertainmentRating: "真实模型判定",
	}
}

func compactJudgeMessage(judgment DebateJudgment) string {
	return fmt.Sprintf("最终裁决：%s 获胜。原因：%s。总结：%s", judgment.WinnerSide, judgment.WinnerReason, judgment.FinalSummary)
}
