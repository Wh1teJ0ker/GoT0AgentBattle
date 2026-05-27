package battle

// DebateStage 表示真实模型辩论中的一个阶段性输出。
// 它既可用于调试，也可作为记忆上下文和裁判材料。
type DebateStage struct {
	Title          string
	Role           string
	SpeakerLabel   string
	SpeakerID      string
	ModelName      string
	ModelProfileID string
	Requirement    string
	Output         string
	MemoryAfter    string
}
