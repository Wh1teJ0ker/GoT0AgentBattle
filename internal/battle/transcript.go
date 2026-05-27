package battle

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func transcriptDir() string {
	return filepath.Join("data", "transcripts")
}

// SaveTranscriptMarkdown 把当前房间的完整对话记录导出为本地 Markdown 文件。
func SaveTranscriptMarkdown(state RoomState) (string, error) {
	if err := os.MkdirAll(transcriptDir(), 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s-%s.md", time.Now().Format("20060102-150405"), sanitizeFilename(state.Topic))
	path := filepath.Join(transcriptDir(), filename)
	content := renderTranscriptMarkdown(state)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return path, nil
	}
	return abs, nil
}

func renderTranscriptMarkdown(state RoomState) string {
	lines := []string{
		"# GoT0AgentBattle 对话记录",
		"",
		fmt.Sprintf("- 主题：`%s`", fallbackText(state.Topic, "未命名辩题")),
		fmt.Sprintf("- 模式：`%s`", fallbackText(state.Mode, "未指定")),
		fmt.Sprintf("- 状态：`%s`", fallbackText(state.Status, "unknown")),
		fmt.Sprintf("- 引擎：`%s`", fallbackText(state.Engine, "unknown")),
		fmt.Sprintf("- 裁判模型档案：`%s`", fallbackText(state.JudgeModelProfileID, "未指定")),
		fmt.Sprintf("- 总回合：`%d`", state.TotalRounds),
		fmt.Sprintf("- 热度：`%d`", state.Heat),
	}

	if state.ProviderPath != "" {
		lines = append(lines, fmt.Sprintf("- 模型配置文件：`%s`", state.ProviderPath))
	}
	if state.StartedAt > 0 {
		lines = append(lines, fmt.Sprintf("- 开始时间：`%s`", formatTimestamp(state.StartedAt)))
	}
	if state.FinishedAt > 0 {
		lines = append(lines, fmt.Sprintf("- 结束时间：`%s`", formatTimestamp(state.FinishedAt)))
	}

	lines = append(lines, "", "## 出场人格", "")
	for _, agent := range state.Agents {
		lines = append(lines,
			fmt.Sprintf("### %s", agent.Name),
			"",
			fmt.Sprintf("- 角色：`%s`", fallbackText(agent.Role, "未指定")),
			fmt.Sprintf("- 风格：`%s`", fallbackText(agent.Style, "未指定")),
			fmt.Sprintf("- 人设：`%s`", fallbackText(agent.Persona, "未指定")),
			fmt.Sprintf("- 模型档案：`%s`", fallbackText(agent.ModelProfileID, "未指定")),
			fmt.Sprintf("- 实际模型：`%s`", fallbackText(agent.Model, "未指定")),
			fmt.Sprintf("- 支持率：`%d%%`", agent.Support),
			fmt.Sprintf("- Token 估算：`%d`", agent.TokenUsage),
			"",
		)
	}

	lines = append(lines, "## 聊天记录", "")
	for _, message := range state.Messages {
		header := fmt.Sprintf("### [%s] %s", formatTimestamp(message.Timestamp), message.AgentName)
		meta := make([]string, 0, 4)
		if message.Round > 0 {
			meta = append(meta, fmt.Sprintf("第 %d 回合", message.Round))
		}
		if message.ReplyTo != "" {
			meta = append(meta, "@"+message.ReplyTo)
		}
		if message.PerformanceTag != "" {
			meta = append(meta, message.PerformanceTag)
		}
		if message.IsJudge {
			meta = append(meta, "裁判")
		}
		if message.IsAudience {
			meta = append(meta, "观众")
		}

		lines = append(lines, header, "")
		if len(meta) > 0 {
			lines = append(lines, "- "+strings.Join(meta, " | "), "")
		}
		lines = append(lines, message.Content, "")
	}

	if state.JudgeSummary != nil {
		lines = append(lines, "## 裁判总结", "")
		lines = append(lines,
			fmt.Sprintf("- 获胜者：`%s`", fallbackText(state.JudgeSummary.WinnerName, "未判定")),
			fmt.Sprintf("- 获胜原因：%s", fallbackText(state.JudgeSummary.WinnerReason, "未提供")),
			fmt.Sprintf("- 节目效果评分：`%d`", state.JudgeSummary.ShowmanshipScore),
			fmt.Sprintf("- 娱乐评级：`%s`", fallbackText(state.JudgeSummary.EntertainmentRating, "未提供")),
			"",
		)

		if len(state.JudgeSummary.EffectiveArguments) > 0 {
			lines = append(lines, "### 有效论据", "")
			for _, item := range state.JudgeSummary.EffectiveArguments {
				lines = append(lines, "- "+item)
			}
			lines = append(lines, "")
		}

		if len(state.JudgeSummary.InvalidTrashTalk) > 0 {
			lines = append(lines, "### 无效嘴炮", "")
			for _, item := range state.JudgeSummary.InvalidTrashTalk {
				lines = append(lines, "- "+item)
			}
			lines = append(lines, "")
		}
	}

	lines = append(lines, "## 排行快照", "")
	for _, agent := range sortAgentsForTranscript(state.Agents) {
		lines = append(lines, fmt.Sprintf("- %s：支持率 `%d%%`，Token `%d`", agent.Name, agent.Support, agent.TokenUsage))
	}
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func sortAgentsForTranscript(agents []AgentState) []AgentState {
	out := append([]AgentState(nil), agents...)
	sort.Slice(out, func(i int, j int) bool {
		if out[i].Support == out[j].Support {
			return out[i].Name < out[j].Name
		}
		return out[i].Support > out[j].Support
	})
	return out
}

func formatTimestamp(ts int64) string {
	if ts <= 0 {
		return "未知时间"
	}
	return time.UnixMilli(ts).Format("2006-01-02 15:04:05")
}

func sanitizeFilename(input string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"\n", "_",
		"\r", "_",
		"\t", "_",
	)
	out := strings.TrimSpace(replacer.Replace(input))
	if out == "" {
		return "debate-transcript"
	}
	return out
}

func fallbackText(input string, fallback string) string {
	if strings.TrimSpace(input) == "" {
		return fallback
	}
	return input
}
