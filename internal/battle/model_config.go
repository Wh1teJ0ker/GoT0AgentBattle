package battle

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ProviderConfig 描述单个模型档案的基础连接信息与重试策略。
type ProviderConfig struct {
	Name              string `json:"name" yaml:"name"`
	BaseURL           string `json:"base_url" yaml:"base_url"`
	APIKey            string `json:"api_key" yaml:"api_key"`
	Model             string `json:"model" yaml:"model"`
	TimeoutSeconds    int    `json:"timeout_seconds" yaml:"timeout_seconds"`
	RetryAttempts     int    `json:"retry_attempts" yaml:"retry_attempts"`
	RetryDelaySeconds int    `json:"retry_delay_seconds" yaml:"retry_delay_seconds"`
}

// DebateProviders 表示可复用的模型档案集合及默认裁判档案。
type DebateProviders struct {
	Profiles              map[string]ProviderConfig `json:"profiles" yaml:"profiles"`
	DefaultJudgeProfileID string                    `json:"default_judge_profile_id" yaml:"default_judge_profile_id"`
}

type legacyDebateProviders struct {
	AffirmativeFirst  ProviderConfig `json:"affirmative_first" yaml:"affirmative_first"`
	AffirmativeSecond ProviderConfig `json:"affirmative_second" yaml:"affirmative_second"`
	AffirmativeThird  ProviderConfig `json:"affirmative_third" yaml:"affirmative_third"`
	NegativeFirst     ProviderConfig `json:"negative_first" yaml:"negative_first"`
	NegativeSecond    ProviderConfig `json:"negative_second" yaml:"negative_second"`
	NegativeThird     ProviderConfig `json:"negative_third" yaml:"negative_third"`
	Judge             ProviderConfig `json:"judge" yaml:"judge"`
}

// GetProfile 根据档案 ID 获取对应的模型配置。
func (p DebateProviders) GetProfile(id string) (ProviderConfig, bool) {
	cfg, ok := p.Profiles[strings.TrimSpace(id)]
	return cfg, ok
}

// Summaries 返回供前端展示使用的模型档案摘要列表。
func (p DebateProviders) Summaries() []ModelProfileSummary {
	keys := make([]string, 0, len(p.Profiles))
	for key := range p.Profiles {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	summaries := make([]ModelProfileSummary, 0, len(keys))
	for _, key := range keys {
		item := p.Profiles[key]
		summaries = append(summaries, ModelProfileSummary{
			ID:          key,
			Label:       item.Name,
			Model:       item.Model,
			Description: strings.TrimSpace(item.BaseURL),
		})
	}
	return summaries
}

// IsReady 判断当前是否至少加载到一个可用模型档案。
func (p DebateProviders) IsReady() bool {
	return len(p.Profiles) > 0
}

func configCandidates() []string {
	candidates := []string{
		os.Getenv("GOT0_AGENT_BATTLE_CONFIG"),
		"got0agentbattle.config.yaml",
		"got0agentbattle.config.yml",
		filepath.Join("config", "got0agentbattle.config.yaml"),
		filepath.Join("config", "got0agentbattle.config.yml"),
		"got0agentbattle.config.json",
		filepath.Join("config", "got0agentbattle.config.json"),
		filepath.Join("config", "providers.json"),
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

// LoadProviderConfig 按候选路径顺序加载模型配置文件。
// 返回值包括解析后的配置、命中的配置文件路径以及错误信息。
func LoadProviderConfig() (DebateProviders, string, error) {
	for _, candidate := range configCandidates() {
		data, err := os.ReadFile(candidate)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return DebateProviders{}, "", err
		}

		cfg, err := parseProviderConfig(data)
		if err != nil {
			return DebateProviders{}, "", fmt.Errorf("invalid provider config %s: %w", candidate, err)
		}
		if err := validateProviderConfig(cfg); err != nil {
			return DebateProviders{}, "", fmt.Errorf("invalid provider config %s: %w", candidate, err)
		}
		abs, err := filepath.Abs(candidate)
		if err != nil {
			return cfg, candidate, nil
		}
		return cfg, abs, nil
	}
	return DebateProviders{}, "", fmt.Errorf("provider config file not found")
}

func parseProviderConfig(data []byte) (DebateProviders, error) {
	var next DebateProviders
	if err := yaml.Unmarshal(data, &next); err == nil && len(next.Profiles) > 0 {
		return next, nil
	}
	if err := json.Unmarshal(data, &next); err == nil && len(next.Profiles) > 0 {
		return next, nil
	}

	var legacy legacyDebateProviders
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		if err := json.Unmarshal(data, &legacy); err != nil {
			return DebateProviders{}, err
		}
	}

	profiles := map[string]ProviderConfig{
		"affirmative-first":  legacy.AffirmativeFirst,
		"affirmative-second": legacy.AffirmativeSecond,
		"affirmative-third":  legacy.AffirmativeThird,
		"negative-first":     legacy.NegativeFirst,
		"negative-second":    legacy.NegativeSecond,
		"negative-third":     legacy.NegativeThird,
		"judge-main":         legacy.Judge,
	}

	aliases := map[string]ProviderConfig{}
	for _, item := range profiles {
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		switch {
		case strings.Contains(strings.ToLower(item.Name), "gpt"), strings.Contains(strings.ToLower(item.Name), "openai"):
			aliases["openai-main"] = item
		case strings.Contains(strings.ToLower(item.Name), "claude"):
			aliases["claude-main"] = item
		case strings.Contains(strings.ToLower(item.Name), "deepseek"):
			aliases["deepseek-main"] = item
		case strings.Contains(strings.ToLower(item.Name), "gemini"), strings.Contains(strings.ToLower(item.Name), "qwen"):
			aliases["gemini-main"] = item
		case strings.Contains(strings.ToLower(item.Name), "ollama"):
			aliases["ollama-local"] = item
		}
	}
	for key, value := range aliases {
		profiles[key] = value
	}

	return DebateProviders{
		Profiles:              profiles,
		DefaultJudgeProfileID: "judge-main",
	}, nil
}

func validateProviderConfig(cfg DebateProviders) error {
	if len(cfg.Profiles) == 0 {
		return fmt.Errorf("missing profiles")
	}

	for id, item := range cfg.Profiles {
		if strings.TrimSpace(item.Name) == "" {
			return fmt.Errorf("missing name for profile: %s", id)
		}
		if strings.TrimSpace(item.BaseURL) == "" {
			return fmt.Errorf("missing base_url for profile: %s", id)
		}
		if strings.TrimSpace(item.APIKey) == "" {
			return fmt.Errorf("missing api_key for profile: %s", id)
		}
		if strings.TrimSpace(item.Model) == "" {
			return fmt.Errorf("missing model for profile: %s", id)
		}
	}

	return nil
}
