package battle

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var randomRoles = []string{
	"云原生嘴替",
	"数据库老炮",
	"性能优化狂",
	"安全研究员",
	"AI 产品经理",
	"前端工程判官",
	"运维值班战神",
	"遗留系统法师",
	"编译器吹毛求疵者",
	"架构评审毒舌位",
}

var randomStyles = []string{
	"阴阳怪气",
	"严谨拆解",
	"冷面补刀",
	"高能拱火",
	"冷嘲短打",
	"梗图连发",
	"火力压制",
	"老狗求生",
}

var randomPersonas = []string{
	"嘴毒但有货",
	"爱抬杠",
	"阴间冷静",
	"直播体质",
	"高压输出",
	"不服就 benchmark",
	"只认事故复盘",
	"梗王本王",
}

var randomTaglines = []string{
	"一句话不一定文明，但一定想让对面闭嘴。",
	"所有争论最后都得上数据，不然就是白吵。",
	"我不负责调和气氛，我负责把问题捅穿。",
	"你可以不服，但你最好先拿出证据。",
	"节目效果和信息增量，我两样都要。",
	"麦一开就默认进入复盘模式。",
}

var randomPromptFocus = []string{
	"优先从工程落地成本、长期维护代价和事故半径切入。",
	"擅长抓逻辑漏洞，逼对方把模糊表述说清楚。",
	"喜欢用数据、基准测试和线上经验狠狠干碎空话。",
	"会主动点名、带节奏，并把观众弹幕变成补刀材料。",
	"遇到不严谨的表达会立刻嘲讽，但不能只有嘴炮没有信息。",
}

var randomColors = []string{
	"#ff6b57",
	"#ff9f43",
	"#4cc9f0",
	"#70a1ff",
	"#2ed573",
	"#ffd166",
	"#a55eea",
	"#ff78c4",
	"#7bed9f",
}

var nonWordPattern = regexp.MustCompile(`[^a-z0-9]+`)

// GenerateRandomPersona 随机生成一个不与现有人格 ID 冲突的新人格模板。
func GenerateRandomPersona(existing []PersonaConfig, judgeModelProfileID string) PersonaConfig {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	role := pick(rng, randomRoles)
	style := pick(rng, randomStyles)
	persona := pick(rng, randomPersonas)
	tagline := pick(rng, randomTaglines)
	color := pick(rng, randomColors)
	modelProfileID := judgeModelProfileID
	if strings.TrimSpace(modelProfileID) == "" {
		modelProfileID = "openai-main"
	}

	nameLeft := []string{"爆表", "断点", "回滚", "火线", "黑盒", "白帽", "高压", "链路", "栈顶", "麦王"}
	nameRight := []string{"阿祖", "判官", "老猫", "暴龙", "修罗", "猎手", "喷子", "战神", "法师", "终端"}
	name := pick(rng, nameLeft) + pick(rng, nameRight)
	id := ensureUniquePersonaID(name, existing)
	avatar := buildAvatar(name)

	return PersonaConfig{
		ID:                    id,
		Name:                  name,
		Role:                  role,
		Avatar:                avatar,
		Color:                 color,
		Style:                 style,
		Persona:               persona,
		Tagline:               tagline,
		Aggressive:            45 + rng.Intn(46),
		Toxicity:              35 + rng.Intn(51),
		Enabled:               true,
		DefaultModelProfileID: modelProfileID,
		SystemPrompt: fmt.Sprintf(
			"你是 %s，一位%s。你的语言风格是%s，性格是%s。%s 你需要在辩论里持续制造压力、抓漏洞、保留信息量，并让节目效果保持在线。",
			name,
			role,
			style,
			persona,
			pick(rng, randomPromptFocus),
		),
	}
}

func ensureUniquePersonaID(name string, existing []PersonaConfig) string {
	base := strings.ToLower(nonWordPattern.ReplaceAllString(name, "-"))
	base = strings.Trim(base, "-")
	if base == "" {
		base = "persona"
	}
	seen := map[string]struct{}{}
	for _, item := range existing {
		seen[item.ID] = struct{}{}
	}
	if _, ok := seen[base]; !ok {
		return base
	}
	for index := 2; ; index++ {
		next := fmt.Sprintf("%s-%d", base, index)
		if _, ok := seen[next]; !ok {
			return next
		}
	}
}

func buildAvatar(name string) string {
	runes := []rune(strings.TrimSpace(name))
	if len(runes) == 0 {
		return "RP"
	}
	if len(runes) == 1 {
		return string(runes[:1])
	}
	return string(runes[:2])
}
