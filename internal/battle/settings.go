package battle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const settingsEnvKey = "GOT0_AGENT_BATTLE_SETTINGS"

var legacyStyleAliases = map[string]string{
	"sarcastic": "阴阳怪气",
	"precision": "严谨拆解",
	"cold":      "冷面补刀",
	"hype":      "高能拱火",
	"deadpan":   "冷嘲短打",
	"meme":      "梗图连发",
	"brutal":    "火力压制",
	"survivor":  "老狗求生",
}

func normalizeStyleLabel(style string) string {
	style = strings.TrimSpace(style)
	if next, ok := legacyStyleAliases[style]; ok {
		return next
	}
	return style
}

func settingsCandidates() []string {
	candidates := []string{
		os.Getenv(settingsEnvKey),
		"got0agentbattle.settings.json",
		filepath.Join("config", "got0agentbattle.settings.json"),
	}

	seen := map[string]struct{}{}
	result := make([]string, 0, len(candidates))
	for _, item := range candidates {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func settingsPath() string {
	candidates := settingsCandidates()
	if len(candidates) == 0 {
		return "got0agentbattle.settings.json"
	}
	return candidates[len(candidates)-1]
}

func defaultPersonas() []PersonaConfig {
	return []PersonaConfig{
		{
			ID: "webdog", Name: "WebDog", Role: "Web 安全", Avatar: "WD", Color: "#ff6b57",
			Style: "阴阳怪气", Persona: "暴躁", Tagline: "喜欢找注入点，也喜欢找人下嘴。",
			Aggressive: 84, Toxicity: 78, Enabled: true, DefaultModelProfileID: "openai-main",
			SystemPrompt: "你是 WebDog，一个暴躁、喜欢点名和讽刺的 Web 安全选手。优先从攻击面、威胁模型、工程落地风险去拆对手观点。",
		},
		{
			ID: "rustman", Name: "RustMan", Role: "Rust 狂热者", Avatar: "RS", Color: "#ff9f43",
			Style: "严谨拆解", Persona: "洁癖", Tagline: "看到未定义行为就想报警。",
			Aggressive: 75, Toxicity: 62, Enabled: true, DefaultModelProfileID: "deepseek-main",
			SystemPrompt: "你是 RustMan，一个对内存安全和类型约束有洁癖的工程师。你会冷静拆解漏洞，再顺手嘲讽不严谨的方案。",
		},
		{
			ID: "opsbuddha", Name: "OpsBuddha", Role: "运维工程师", Avatar: "OP", Color: "#2ed573",
			Style: "冷面补刀", Persona: "克制", Tagline: "凡是不能过夜班的方案都不值得吹。",
			Aggressive: 58, Toxicity: 40, Enabled: true, DefaultModelProfileID: "claude-main",
			SystemPrompt: "你是 OpsBuddha，一个对线上稳定性极度敏感的运维工程师。你优先讨论 SLO、回滚、值班代价和事故半径。",
		},
		{
			ID: "promptmax", Name: "PromptMax", Role: "AI 狂热者", Avatar: "AI", Color: "#4cc9f0",
			Style: "高能拱火", Persona: "上头", Tagline: "上下文一拉满，谁都敢点评。",
			Aggressive: 70, Toxicity: 55, Enabled: true, DefaultModelProfileID: "gemini-main",
			SystemPrompt: "你是 PromptMax，一个热衷把任何话题往 AI 自动化和模型能力上拽的人。你说话节奏快，爱放大趋势和组织效率差异。",
		},
		{
			ID: "binarycat", Name: "BinaryCat", Role: "二进制选手", Avatar: "BC", Color: "#a55eea",
			Style: "冷嘲短打", Persona: "阴阳怪气", Tagline: "一句 UAF 能顶你三页 PPT。",
			Aggressive: 88, Toxicity: 83, Enabled: true, DefaultModelProfileID: "ollama-local",
			SystemPrompt: "你是 BinaryCat，一个阴阳怪气、擅长边界条件和底层异常路径攻击的二进制选手。你偏爱短句、冷嘲和杀伤力。",
		},
		{
			ID: "jspirate", Name: "JsPirate", Role: "前端海盗", Avatar: "JS", Color: "#ffd166",
			Style: "梗图连发", Persona: "玩梗", Tagline: "npm 装一切，也敢喷一切。",
			Aggressive: 66, Toxicity: 69, Enabled: true, DefaultModelProfileID: "openai-main",
			SystemPrompt: "你是 JsPirate，一个满嘴梗、擅长从生态复杂度和开发体验角度点燃气氛的前端海盗。",
		},
		{
			ID: "kernelbro", Name: "KernelBro", Role: "系统底层党", Avatar: "KB", Color: "#7bed9f",
			Style: "火力压制", Persona: "硬核", Tagline: "性能不够就是原罪，解释权归火焰图所有。",
			Aggressive: 80, Toxicity: 60, Enabled: true, DefaultModelProfileID: "deepseek-main",
			SystemPrompt: "你是 KernelBro，一个执着于吞吐、延迟、资源占用和火焰图的系统底层党。你会用性能和成本数据狠狠干碎空话。",
		},
		{
			ID: "phpghost", Name: "PhpGhost", Role: "遗留系统守护者", Avatar: "PG", Color: "#70a1ff",
			Style: "老狗求生", Persona: "老油条", Tagline: "你们都在画未来，我负责把线上活到明天。",
			Aggressive: 61, Toxicity: 58, Enabled: true, DefaultModelProfileID: "claude-main",
			SystemPrompt: "你是 PhpGhost，一个专门从遗留系统、迁移成本、业务连续性和现实妥协角度拆台的人。",
		},
	}
}

func defaultSettings() AppSettings {
	cfg := defaultConfig()
	cfg.PersonaIDs = []string{"webdog", "rustman", "opsbuddha", "promptmax", "binarycat"}
	cfg.JudgeModelProfileID = "judge-main"
	cfg.PreferRealLLM = true
	return AppSettings{
		DefaultConfig:        cfg,
		Personas:             defaultPersonas(),
		JudgeModelProfileID:  "judge-main",
		PreferredThemeFlavor: "broadcast-chaos",
	}
}

func normalizeSettings(settings AppSettings) AppSettings {
	out := settings
	if len(out.Personas) == 0 {
		out.Personas = defaultPersonas()
	}
	for i := range out.Personas {
		persona := &out.Personas[i]
		persona.ID = strings.TrimSpace(persona.ID)
		if persona.ID == "" {
			persona.ID = fmt.Sprintf("persona-%d", i+1)
		}
		if strings.TrimSpace(persona.Name) == "" {
			persona.Name = "未命名人格"
		}
		if strings.TrimSpace(persona.Role) == "" {
			persona.Role = "自由辩手"
		}
		if strings.TrimSpace(persona.Avatar) == "" {
			persona.Avatar = "AG"
		}
		if strings.TrimSpace(persona.Color) == "" {
			persona.Color = "#8b95aa"
		}
		persona.Style = normalizeStyleLabel(persona.Style)
		if strings.TrimSpace(persona.Style) == "" {
			persona.Style = "阴阳怪气"
		}
		if strings.TrimSpace(persona.Persona) == "" {
			persona.Persona = "嘴硬"
		}
		if strings.TrimSpace(persona.Tagline) == "" {
			persona.Tagline = "正在寻找下一句更狠的话。"
		}
		if persona.Aggressive <= 0 {
			persona.Aggressive = 60
		}
		if persona.Toxicity <= 0 {
			persona.Toxicity = 50
		}
		if strings.TrimSpace(persona.DefaultModelProfileID) == "" {
			persona.DefaultModelProfileID = "openai-main"
		}
		if strings.TrimSpace(persona.SystemPrompt) == "" {
			persona.SystemPrompt = fmt.Sprintf("你是 %s，一个%s风格的%s。你会围绕工程现实、观点漏洞和对手弱点展开火力。", persona.Name, persona.Persona, persona.Role)
		}
	}

	out.DefaultConfig = normalizeConfig(out.DefaultConfig)
	if len(out.DefaultConfig.PersonaIDs) == 0 {
		for _, persona := range out.Personas {
			if persona.Enabled {
				out.DefaultConfig.PersonaIDs = append(out.DefaultConfig.PersonaIDs, persona.ID)
			}
			if len(out.DefaultConfig.PersonaIDs) >= out.DefaultConfig.AgentCount {
				break
			}
		}
	}
	if strings.TrimSpace(out.JudgeModelProfileID) == "" {
		out.JudgeModelProfileID = "judge-main"
	}
	if strings.TrimSpace(out.DefaultConfig.JudgeModelProfileID) == "" {
		out.DefaultConfig.JudgeModelProfileID = out.JudgeModelProfileID
	}
	if strings.TrimSpace(out.PreferredThemeFlavor) == "" {
		out.PreferredThemeFlavor = "broadcast-chaos"
	}
	return out
}

// LoadSettings 从本地候选路径中加载设置文件。
// 如果没有找到文件，则返回默认设置和推荐保存路径。
func LoadSettings() (AppSettings, string, error) {
	for _, candidate := range settingsCandidates() {
		data, err := os.ReadFile(candidate)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return AppSettings{}, "", err
		}

		var settings AppSettings
		if err := json.Unmarshal(data, &settings); err != nil {
			return AppSettings{}, "", fmt.Errorf("invalid settings %s: %w", candidate, err)
		}

		abs, err := filepath.Abs(candidate)
		if err != nil {
			return normalizeSettings(settings), candidate, nil
		}
		return normalizeSettings(settings), abs, nil
	}

	path := settingsPath()
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return defaultSettings(), abs, nil
}

// SaveSettings 把设置写入本地文件，并返回归一化后的结果与实际保存路径。
func SaveSettings(settings AppSettings) (AppSettings, string, error) {
	normalized := normalizeSettings(settings)
	path := settingsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return AppSettings{}, "", err
	}

	payload, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		return AppSettings{}, "", err
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return AppSettings{}, "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return normalized, path, nil
	}
	return normalized, abs, nil
}

func cloneSettings(settings AppSettings) AppSettings {
	copySettings := settings
	copySettings.DefaultConfig.PersonaIDs = append([]string(nil), settings.DefaultConfig.PersonaIDs...)
	copySettings.Personas = append([]PersonaConfig(nil), settings.Personas...)
	return copySettings
}
