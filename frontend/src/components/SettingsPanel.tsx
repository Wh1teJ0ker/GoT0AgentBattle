// 文件说明：SettingsPanel 负责承载独立设置中心页面。
// 职责：维护默认房间参数、本地人格库、主题偏好以及设置持久化操作入口。
import {
  Alert,
  Button,
  Card,
  Checkbox,
  Col,
  Divider,
  Flex,
  Form,
  Input,
  InputNumber,
  Row,
  Select,
  Space,
  Switch,
  Typography,
} from 'antd';
import type {AppSettings, BootstrapData, PersonaConfig, RoomForm} from '../types';

const {Paragraph, Text, Title} = Typography;
const {TextArea} = Input;

type SettingsPanelProps = {
  boot: BootstrapData;
  settings: AppSettings;
  setSettings: (next: AppSettings) => void;
  busy: boolean;
  saveBusy: boolean;
  settingsPath: string;
  onSave: () => Promise<void>;
  onGenerateRandomPersona: () => Promise<void>;
};

function SettingsPanel({
  boot,
  settings,
  setSettings,
  busy,
  saveBusy,
  settingsPath,
  onSave,
  onGenerateRandomPersona,
}: SettingsPanelProps) {
  const selectedPersonaIDs = settings.defaultConfig.personaIds ?? [];

  function updateDefaultConfig(next: RoomForm) {
    setSettings({...settings, defaultConfig: next});
  }

  function updatePersona(index: number, patch: Partial<PersonaConfig>) {
    const personas = settings.personas.map((item, itemIndex) => itemIndex === index ? {...item, ...patch} : item);
    setSettings({...settings, personas});
  }

  function togglePersonaSelection(id: string, checked: boolean) {
    const nextSet = new Set(selectedPersonaIDs);
    if (checked) {
      nextSet.add(id);
    } else {
      nextSet.delete(id);
    }
    updateDefaultConfig({
      ...settings.defaultConfig,
      personaIds: Array.from(nextSet),
      agentCount: Math.max(1, Array.from(nextSet).length),
    });
  }

  return (
    <div className="settings-shell">
      <Row gutter={[18, 18]}>
        <Col xs={24} xl={9}>
          <Space direction="vertical" size={18} className="full-width">
            <Card className="battle-panel settings-panel" bordered={false}>
              <Text className="panel-kicker">设置中心</Text>
              <Title level={2}>总控台</Title>
              <Paragraph>在这里配置默认房间、前端风格和人格库。模型提供方、档案与裁判默认模型统一由 YAML 配置文件控制。</Paragraph>
              <Alert
                type="info"
                showIcon
                message="本地设置文件路径"
                description={settingsPath || 'got0agentbattle.settings.json'}
              />
              <Alert
                type="warning"
                showIcon
                message="模型设置已迁出设置中心"
                description={boot.realProviderPath || '请在 got0agentbattle.config.yaml 中维护模型档案、API Key 和默认裁判模型。'}
              />
              <Divider/>
              <Form layout="vertical">
                <Form.Item label="默认辩论模式">
                  <Select
                    value={settings.defaultConfig.mode}
                    onChange={(value) => updateDefaultConfig({...settings.defaultConfig, mode: value})}
                    options={boot.modes.map((mode) => ({value: mode.id, label: mode.label}))}
                  />
                </Form.Item>
                <Row gutter={12}>
                  <Col span={12}>
                    <Form.Item label="默认回合数">
                      <InputNumber
                        min={1}
                        max={6}
                        className="full-width"
                        value={settings.defaultConfig.rounds}
                        onChange={(value) => updateDefaultConfig({...settings.defaultConfig, rounds: Number(value || 3)})}
                      />
                    </Form.Item>
                  </Col>
                  <Col span={12}>
                    <Form.Item label="默认上场人数">
                      <InputNumber
                        min={1}
                        max={12}
                        className="full-width"
                        value={settings.defaultConfig.agentCount}
                        onChange={(value) => updateDefaultConfig({...settings.defaultConfig, agentCount: Number(value || selectedPersonaIDs.length || 1)})}
                      />
                    </Form.Item>
                  </Col>
                </Row>
                <Form.Item label="主题风格">
                  <Select
                    value={settings.preferredThemeFlavor}
                    onChange={(value) => setSettings({...settings, preferredThemeFlavor: value})}
                    options={[
                      {value: 'broadcast-chaos', label: '直播混战'},
                      {value: 'arcade-fever', label: '街机热区'},
                      {value: 'iron-war-room', label: '铁血导播间'},
                    ]}
                  />
                </Form.Item>
                <Button type="primary" size="large" onClick={() => void onSave()} loading={saveBusy} disabled={busy} block>
                  保存设置
                </Button>
              </Form>
            </Card>
          </Space>
        </Col>

        <Col xs={24} xl={15}>
          <Card className="battle-panel settings-panel" bordered={false}>
            <Flex justify="space-between" align="center" wrap="wrap" gap={12}>
              <div>
                <Text className="panel-kicker">人格库</Text>
                <Title level={2}>角色导演台</Title>
              </div>
              <Space>
                <Button onClick={() => void onGenerateRandomPersona()}>随机生成人格</Button>
              </Space>
            </Flex>
            <Paragraph>勾选上场名单，维护人设、风格和系统提示词。模型档案绑定保留在本地数据结构中，但不再在这里直接编辑。</Paragraph>
            <Space direction="vertical" size={16} className="full-width">
              {settings.personas.map((persona, index) => (
                <Card key={persona.id} className="persona-editor" bordered={false}>
                  <Row gutter={[12, 12]}>
                    <Col xs={24} md={8}>
                      <Form layout="vertical">
                        <Form.Item label="上场">
                          <Checkbox
                            checked={selectedPersonaIDs.includes(persona.id)}
                            onChange={(event) => togglePersonaSelection(persona.id, event.target.checked)}
                          >
                            作为默认出战人格
                          </Checkbox>
                        </Form.Item>
                        <Form.Item label="启用人格">
                          <Switch
                            checked={persona.enabled}
                            onChange={(checked: boolean) => updatePersona(index, {enabled: checked})}
                          />
                        </Form.Item>
                        <Form.Item label="名字">
                          <Input value={persona.name} onChange={(event) => updatePersona(index, {name: event.target.value})}/>
                        </Form.Item>
                        <Form.Item label="头像字母">
                          <Input value={persona.avatar} onChange={(event) => updatePersona(index, {avatar: event.target.value.slice(0, 3)})}/>
                        </Form.Item>
                        <Form.Item label="颜色">
                          <Input value={persona.color} onChange={(event) => updatePersona(index, {color: event.target.value})}/>
                        </Form.Item>
                      </Form>
                    </Col>
                    <Col xs={24} md={8}>
                      <Form layout="vertical">
                        <Form.Item label="专业方向">
                          <Input value={persona.role} onChange={(event) => updatePersona(index, {role: event.target.value})}/>
                        </Form.Item>
                        <Form.Item label="性格标签">
                          <Input value={persona.persona} onChange={(event) => updatePersona(index, {persona: event.target.value})}/>
                        </Form.Item>
                        <Form.Item label="语言风格">
                          <Select
                            value={persona.style}
                            onChange={(value) => updatePersona(index, {style: value})}
                            options={[
                              {value: '阴阳怪气', label: '阴阳怪气'},
                              {value: '严谨拆解', label: '严谨拆解'},
                              {value: '冷面补刀', label: '冷面补刀'},
                              {value: '高能拱火', label: '高能拱火'},
                              {value: '冷嘲短打', label: '冷嘲短打'},
                              {value: '梗图连发', label: '梗图连发'},
                              {value: '火力压制', label: '火力压制'},
                              {value: '老狗求生', label: '老狗求生'},
                            ]}
                          />
                        </Form.Item>
                        <Form.Item label="当前模型绑定">
                          <Input value={persona.defaultModelProfileId} disabled/>
                        </Form.Item>
                        <Row gutter={12}>
                          <Col span={12}>
                            <Form.Item label="攻击性">
                              <InputNumber
                                min={0}
                                max={100}
                                className="full-width"
                                value={persona.aggressive}
                                onChange={(value) => updatePersona(index, {aggressive: Number(value || 0)})}
                              />
                            </Form.Item>
                          </Col>
                          <Col span={12}>
                            <Form.Item label="嘴臭程度">
                              <InputNumber
                                min={0}
                                max={100}
                                className="full-width"
                                value={persona.toxicity}
                                onChange={(value) => updatePersona(index, {toxicity: Number(value || 0)})}
                              />
                            </Form.Item>
                          </Col>
                        </Row>
                      </Form>
                    </Col>
                    <Col xs={24} md={8}>
                      <Form layout="vertical">
                        <Form.Item label="角色口号">
                          <Input value={persona.tagline} onChange={(event) => updatePersona(index, {tagline: event.target.value})}/>
                        </Form.Item>
                        <Form.Item label="系统 Prompt">
                          <TextArea
                            value={persona.systemPrompt}
                            onChange={(event) => updatePersona(index, {systemPrompt: event.target.value})}
                            autoSize={{minRows: 8, maxRows: 14}}
                          />
                        </Form.Item>
                      </Form>
                    </Col>
                  </Row>
                </Card>
              ))}
            </Space>
          </Card>
        </Col>
      </Row>
    </div>
  );
}

export default SettingsPanel
