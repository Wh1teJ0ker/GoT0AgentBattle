package battle

import (
	"fmt"
	"strings"
)

func buildPersonaPrompt(topic string, mode string, round int, speaker AgentState, target AgentState, memory string, action *AudienceActionInput, persona PersonaConfig) string {
	actionHint := "无额外观众干预。"
	actionRule := "按你自己的节奏打，但要制造可被下一位接上的转火点。"
	if action != nil {
		actionHint = fmt.Sprintf("观众刚刚插话：[%s] %s。", action.Kind, action.Message)
		if strings.TrimSpace(action.TargetAgentID) != "" {
			actionRule = "这条场外指令带了指定攻击对象，你必须顺势把火力集中到该对象。"
		} else {
			actionRule = "你要把这条场外节奏自然接进自己的发言里，像现场接梗，而不是生硬引用。"
		}
	}

	return fmt.Sprintf(`你正在参与一个节目化、娱乐化的 AI 多人辩论秀。

当前辩题：
%s

辩论模式：
%s

当前回合：
第 %d 回合

你的人设资料：
- 名称：%s
- 专业方向：%s
- 性格：%s
- 语言风格：%s
- 攻击性：%d / 100
- 嘴臭程度：%d / 100
- 台词设定：%s

额外角色设定：
%s

你这回合的主要攻击对象：
- 名称：%s
- 方向：%s
- 最近一句话：%s

历史发言摘要：
%s

导播信息：
%s

当前舞台态势：
- 你的当前愤怒值：%d / 100
- 你的当前支持率：%d / 100
- 你的当前 momentum：%d / 100
- 对方上一轮针对的人：%s

你的任务：
1. 必须围绕辩题和对手弱点发言，不要跑题
2. 至少点名一次 @%s
3. 要有明确论点、反驳和嘲讽，但不能只有空洞脏话
4. 允许有节目效果、梗和攻击性，但要保留工程信息量
5. 最好顺带制造转火点，方便下一位继续互喷
6. %s

输出要求：
- 直接输出一段中文发言
- 2 到 6 句
 - 不要写标题，不要解释自己在做什么`, topic, mode, round, speaker.Name, speaker.Role, speaker.Persona, speaker.Style, speaker.Aggressive, speaker.Toxicity, speaker.Tagline, persona.SystemPrompt, target.Name, target.Role, fallbackLine(target.LastLine), memory, actionHint, speaker.Anger, speaker.Support, speaker.Momentum, fallbackLine(target.LastTarget), target.Name, actionRule)
}

func buildJudgePrompt(topic string, round int, memory string, agents []AgentState) string {
	roster := make([]string, 0, len(agents))
	for _, agent := range agents {
		roster = append(roster, fmt.Sprintf("- %s（%s, model=%s）", agent.Name, agent.Role, agent.Model))
	}

	return fmt.Sprintf(`你是这档 AI 辩论节目的裁判 Judge.exe。

辩题：
%s

当前轮次：
第 %d 回合

本场选手：
%s

完整辩论记录：
%s

请你完成裁决，要求兼顾“信息量”和“节目效果”。

请只输出 JSON：
{
  "winner_side": "这里写获胜者名字",
  "winner_reason": "为什么他/她/它赢了",
  "speaker_scores": [
    {
      "speaker": "选手名字",
      "model": "使用的模型",
      "breakdown": {
        "argument": 0,
        "evidence": 0,
        "rebuttal": 0,
        "clarity": 0,
        "strategy": 0,
        "total": 0
      },
      "comment": "简短评价"
    }
  ],
  "side_totals": {
    "showmanship": 0,
    "substance": 0
  },
  "key_clashes": [
    "关键交锋 1",
    "关键交锋 2",
    "关键交锋 3"
  ],
  "final_summary": "对本回合的整体总结"
}

要求：
- winner_side 直接写获胜人格名，不要写 affirmative/negative
- speaker_scores 至少覆盖全部出场选手
- total 等于前五项之和
- final_summary 要像节目裁判发言，不要像论文`, topic, round, strings.Join(roster, "\n"), memory)
}

func judgeRequirement() string {
	return "结合节目效果与论证质量，为所有出场选手打分，选出获胜人格，给出关键交锋与总结。"
}

func fallbackLine(line string) string {
	if strings.TrimSpace(line) == "" {
		return "暂时还没开口。"
	}
	return line
}
