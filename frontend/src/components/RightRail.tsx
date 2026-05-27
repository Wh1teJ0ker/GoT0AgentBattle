// 文件说明：RightRail 负责承载战斗页右侧的房间控制与信息面板。
// 职责：管理建房表单、状态提示、排行榜摘要与开战操作入口。
import {
  Alert,
  Button,
  Card,
  Col,
  Divider,
  Flex,
  Form,
  Input,
  InputNumber,
  List,
  Row,
  Select,
  Space,
  Statistic,
  Typography,
} from 'antd';
import type {BootstrapData, RoomForm, RoomState} from '../types';

const {Paragraph, Text, Title} = Typography;
const {TextArea} = Input;

type RightRailProps = {
  boot: BootstrapData;
  state: RoomState;
  form: RoomForm;
  setForm: (next: RoomForm) => void;
  busy: boolean;
  error: string;
  topPerformer?: RoomState['agents'][number];
  supportRanking: RoomState['agents'];
  onCreateRoom: () => Promise<void>;
  onStartBattle: () => Promise<void>;
};

function RightRail({
  boot,
  state,
  form,
  setForm,
  busy,
  error,
  topPerformer,
  supportRanking,
  onCreateRoom,
  onStartBattle,
}: RightRailProps) {
  return (
    <Space direction="vertical" size={18} className="full-width">
      <Card className="battle-panel" bordered={false}>
        <Text className="panel-kicker">导播配置</Text>
        <Title level={3}>房间控制台</Title>
        <Form layout="vertical">
          <Form.Item label="辩论主题">
            <TextArea
              value={form.topic}
              onChange={(event) => setForm({...form, topic: event.target.value})}
              placeholder="手动输入本场辩题，例如：Go 和 Rust 谁更适合后端开发？"
              autoSize={{minRows: 3, maxRows: 5}}
            />
          </Form.Item>
          <Divider/>
          <Form.Item label="辩论模式">
            <Select
              value={form.mode}
              onChange={(value) => setForm({...form, mode: value})}
              options={boot.modes.map((mode) => ({value: mode.id, label: mode.label}))}
            />
          </Form.Item>
          <Form.Item label="默认上场人格">
            <Select
              mode="multiple"
              value={form.personaIds}
              onChange={(value) => setForm({...form, personaIds: value, agentCount: value.length || form.agentCount})}
              options={boot.settings.personas
                .filter((persona) => persona.enabled)
                .map((persona) => ({value: persona.id, label: persona.name}))}
            />
          </Form.Item>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item label="角色数量">
                <InputNumber
                  min={3}
                  max={8}
                  value={form.agentCount}
                  onChange={(value) => setForm({...form, agentCount: Number(value || 5)})}
                  className="full-width"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item label="回合数">
                <InputNumber
                  min={1}
                  max={6}
                  value={form.rounds}
                  onChange={(value) => setForm({...form, rounds: Number(value || 3)})}
                  className="full-width"
                />
              </Form.Item>
            </Col>
          </Row>
          <Alert
            type={state.providerReady ? 'success' : 'warning'}
            showIcon
            message={state.providerNotice}
            description={state.providerPath || boot.realProviderPath || '未检测到 got0agentbattle.config.yaml。'}
          />
          <Alert
            type="info"
            showIcon
            message="模型切换已改为配置文件驱动"
            description="请直接修改本地 YAML 配置文件来维护模型档案、API Key 与默认裁判模型。"
          />
          <Divider/>
          <Flex gap={10}>
            <Button type="primary" onClick={() => void onCreateRoom()} disabled={busy} block>
              创建房间
            </Button>
            <Button danger onClick={() => void onStartBattle()} disabled={busy || !state.roomId || !state.providerReady} block>
              开始互喷
            </Button>
          </Flex>
        </Form>
      </Card>

      <Card className="battle-panel" bordered={false}>
        <Text className="panel-kicker">实时排行</Text>
        <Title level={4}>{topPerformer ? `${topPerformer.name} 领跑中` : '等待数据'}</Title>
        <Paragraph>{state.currentTarget ? `当前被集中围攻：${state.currentTarget}` : '还没形成明显集火对象。'}</Paragraph>
        <List
          dataSource={supportRanking}
          renderItem={(agent, index) => (
            <List.Item>
              <Flex justify="space-between" align="center" className="full-width">
                <Space>
                  <Text className="rank-index">#{index + 1}</Text>
                  <Text strong>{agent.name}</Text>
                </Space>
                <Space size={10}>
                  <Text type="secondary">势能 {agent.momentum}</Text>
                  <Text>{agent.support}%</Text>
                </Space>
              </Flex>
            </List.Item>
          )}
        />
      </Card>

      <Card className="battle-panel battle-panel-judge" bordered={false}>
        <Text className="panel-kicker">裁判结论</Text>
        <Title level={4}>{state.judgeSummary ? `${state.judgeSummary.winnerName} 拿下本轮` : '尚未判决'}</Title>
        {state.judgeSummary ? (
          <Space direction="vertical" size={12} className="full-width">
            <Paragraph>{state.judgeSummary.winnerReason}</Paragraph>
            <Statistic title="节目效果评分" value={state.judgeSummary.showmanshipScore}/>
            <div>
              {state.judgeSummary.effectiveArguments.map((item) => (
                <Paragraph key={item} className="judge-point">
                  + {item}
                </Paragraph>
              ))}
            </div>
            <div>
              {state.judgeSummary.invalidTrashTalk.map((item) => (
                <Paragraph key={item} className="judge-point judge-point-negative">
                  - {item}
                </Paragraph>
              ))}
            </div>
          </Space>
        ) : (
          <Paragraph>每回合结束后，Judge.exe 会在这里宣布胜者与节目效果评分。</Paragraph>
        )}
      </Card>

      {error ? <Alert type="error" showIcon message={error}/> : null}
    </Space>
  );
}

export default RightRail
