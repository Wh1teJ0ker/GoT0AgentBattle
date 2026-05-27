package battle

import (
	"encoding/json"
	"strings"
)

// ScoreBreakdown 表示裁判对单个发言者的各维度拆分得分。
type ScoreBreakdown struct {
	Argument int `json:"argument"`
	Evidence int `json:"evidence"`
	Rebuttal int `json:"rebuttal"`
	Clarity  int `json:"clarity"`
	Strategy int `json:"strategy"`
	Total    int `json:"total"`
}

// SpeakerScore 表示裁判对某个发言者的评分结果。
type SpeakerScore struct {
	Speaker   string         `json:"speaker"`
	Model     string         `json:"model"`
	Breakdown ScoreBreakdown `json:"breakdown"`
	Comment   string         `json:"comment"`
}

// DebateJudgment 表示裁判模型输出的整场辩论裁决结构。
type DebateJudgment struct {
	WinnerSide    string         `json:"winner_side"`
	WinnerReason  string         `json:"winner_reason"`
	SpeakerScores []SpeakerScore `json:"speaker_scores"`
	SideTotals    map[string]int `json:"side_totals"`
	KeyClashes    []string       `json:"key_clashes"`
	FinalSummary  string         `json:"final_summary"`
}

func parseJudgment(raw string) (DebateJudgment, error) {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var out DebateJudgment
	err := json.Unmarshal([]byte(cleaned), &out)
	return out, err
}
