package battle

import (
	"fmt"
	"math/rand"
	"strings"
)

var styleOpeners = map[string][]string{
	"阴阳怪气": {
		"你先别急着赢，先把逻辑缝起来。",
		"这套说辞听起来很满，细看全是洞。",
	},
	"严谨拆解": {
		"我先校验一下你这套论证的内存安全。",
		"结论先别飞，约束条件还没落地。",
	},
	"冷面补刀": {
		"我只看能不能过凌晨三点的值班。",
		"别喊口号，先算事故半径。",
	},
	"高能拱火": {
		"我承认你很会说，但上下文窗口已经给你判死刑了。",
		"这题我有感觉，接下来是 token 级别的转火。",
	},
	"冷嘲短打": {
		"你这发言像野指针，碰巧没炸不代表安全。",
		"我听完只有一个反应：这玩意怎么还没 core dump。",
	},
	"梗图连发": {
		"这逻辑我建议先 npm unpublish，再重新做人。",
		"你说得像是能上线，结果一跑全是梗图。",
	},
	"火力压制": {
		"性能图一摊开，嘴硬的人都会安静。",
		"别扯抽象叙事，火焰图会替我说话。",
	},
	"老狗求生": {
		"你们讨论理想世界，我负责看明早谁背锅。",
		"线上系统没空听神话，它只认回滚按钮。",
	},
}

var roleAngles = map[string][]string{
	"Web 安全": {
		"没有威胁建模的方案，最后都得靠补锅来买单",
		"只要攻击面没收缩，再漂亮的架构都只是更高级的靶子",
	},
	"Rust 狂热者": {
		"工程效率不能靠把未定义行为当人品测试来换",
		"后端稳定性不是情怀，是每次发布之后还能不能睡觉",
	},
	"运维工程师": {
		"上线之后的可观测性和回滚路径，决定这东西有没有资格谈价值",
		"任何不考虑 SLO 和故障恢复的观点，都是没见过凌晨报警",
	},
	"AI 狂热者": {
		"自动化放大的不仅是效率，还会放大你原本的组织水平",
		"模型能不能打，不只看智商，还看成本和上下文纪律",
	},
	"二进制选手": {
		"真正能打的方案，必须经得住边界条件和异常路径的拷打",
		"别把运行成功当安全，只要一崩就是一整片事故现场",
	},
	"前端海盗": {
		"复杂度一旦堆起来，最终受苦的是每个接手的人",
		"生态再热闹，不能维护就是给未来的自己埋雷",
	},
	"系统底层党": {
		"吞吐、延迟和资源占用才是最终裁判",
		"没有性能余量的优雅，线上根本活不到第二天",
	},
	"遗留系统守护者": {
		"真正的成熟不是新，而是出事时还能稳住业务",
		"别笑旧系统，笑到最后的往往是那个最难挂的家伙",
	},
}

var attackClosers = []string{
	"你现在不是在论证，你是在拿气势冒充结论。",
	"节目效果可以有，但论据不能靠音量生成。",
	"拿情绪补论证，是这场最偷懒的写法。",
	"如果这也算方案，那 PPT 模板也能当架构图了。",
}

var actionResponses = []string{
	"弹幕这刀补得挺准，我顺手就接了。",
	"场外观众这句提得对，刚好暴露了你的软肋。",
	"谢谢弹幕递话筒，我就按这个角度继续拆。",
}

var hypeTags = []string{
	"开始转火",
	"当场拆台",
	"火力覆盖",
	"情绪拉满",
	"逻辑爆破",
}

func composeAgentLine(rng *rand.Rand, topic string, mode string, round int, speaker AgentState, target AgentState, action *AudienceActionInput) (string, string, string) {
	opener := pick(rng, styleOpeners[speaker.Style])
	angle := pick(rng, roleAngles[speaker.Role])
	targetPunch := fmt.Sprintf("@%s 你最大的问题是把“%s”讲成了立场表演。", target.Name, topic)
	modeLine := modeFlavor(mode)
	closeLine := pick(rng, attackClosers)
	parts := []string{
		opener,
		fmt.Sprintf("拿“%s”这题来说，%s。", topic, angle),
		targetPunch,
		modeLine,
	}
	if action != nil {
		parts = append(parts, fmt.Sprintf("%s 观众刚刷“%s”，这句我直接拿来当补刀。", pick(rng, actionResponses), action.Message))
	}
	if round >= 2 {
		parts = append(parts, "到了这个回合还在兜圈子，只能说明你没有新的有效信息。")
	}
	parts = append(parts, closeLine)
	return strings.Join(parts, " "), pick(rng, hypeTags), speaker.Style
}

func focusByRole(role string) string {
	items := roleAngles[role]
	if len(items) == 0 {
		return "真实工程代价"
	}
	return items[0]
}

func modeFlavor(mode string) string {
	switch mode {
	case "free-for-all":
		return "自由辩论不是自由漂移，谁偏题谁先掉分。"
	case "red-vs-blue":
		return "既然是对抗局，就别拿模糊表态给自己留后门。"
	case "host-mode":
		return "主持人还没打断你，说明我现在拆得还不够狠。"
	case "king-of-hill":
		return "擂台上只认站到最后的人，不认发言最长的人。"
	default:
		return "混战模式最怕的不是吵，而是有人吵半天还没有实质输出。"
	}
}

func pick(rng *rand.Rand, items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[rng.Intn(len(items))]
}
