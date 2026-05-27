package battle

import (
	"strconv"
	"strings"
)

// DebateMemory 保存真实模型辩论过程中的上下文记忆与阶段日志。
type DebateMemory struct {
	topic  string
	stages []DebateStage
}

// NewMemory 创建一份用于真实模型辩论过程的上下文记忆。
func NewMemory(topic string) *DebateMemory {
	return &DebateMemory{
		topic: topic,
	}
}

// Add 追加一个新的辩论阶段到记忆中。
func (m *DebateMemory) Add(stage DebateStage) {
	m.stages = append(m.stages, stage)
}

// Snapshot 返回适合继续喂给模型的简化上下文摘要。
func (m *DebateMemory) Snapshot() string {
	var parts []string
	parts = append(parts, "Topic: "+m.topic)
	for i, stage := range m.stages {
		parts = append(parts, "Stage "+strconv.Itoa(i+1)+" - "+stage.SpeakerLabel+" ("+stage.ModelName+"):")
		parts = append(parts, stage.Output)
	}
	return strings.Join(parts, "\n")
}

// FullLog 返回包含完整阶段信息的上下文日志，适合裁判总结阶段使用。
func (m *DebateMemory) FullLog() string {
	var parts []string
	parts = append(parts, "Full debate memory log:")
	for i, stage := range m.stages {
		parts = append(parts, "")
		parts = append(parts, "Stage "+strconv.Itoa(i+1)+": "+stage.Title)
		parts = append(parts, "Speaker: "+stage.SpeakerLabel)
		parts = append(parts, "Model: "+stage.ModelName)
		parts = append(parts, "Requirement: "+stage.Requirement)
		parts = append(parts, "Output:")
		parts = append(parts, stage.Output)
	}
	return strings.Join(parts, "\n")
}
