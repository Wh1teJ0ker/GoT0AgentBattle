// 文件说明：BattleStage 负责承载主战场界面。
// 职责：组合角色席位、聊天室、舞台指标、观众互动区与右侧控制面板。
import {
  Alert,
  Badge,
  Button,
  Card,
  Col,
  Flex,
  Input,
  Layout,
  Progress,
  Row,
  Select,
  Space,
  Tag,
  Typography,
} from 'antd';
import {Suspense} from 'react';
import type {BootstrapData, RoomForm, RoomState} from '../types';
import AgentRosterPanel from './AgentRosterPanel';
import RightRail from './RightRail';
import ChatFeed from './ChatFeed';
import {MetricCard, PanelFallback, modeLabel, stagePhaseLabel, statusLabel} from '../lib/battle-ui';

const {Sider, Content} = Layout;
const {Title, Paragraph, Text} = Typography;
const {TextArea} = Input;

const quickActions = [
  {label: '攻击 PHP', kind: '指令'},
  {label: '讨论内存安全', kind: '弹幕'},
  {label: '别吵了，上证据', kind: '主持人'},
  {label: '转火性能瓶颈', kind: '导播'},
];

type BattleStageProps = {
  boot: BootstrapData;
  state: RoomState;
  form: RoomForm;
  error: string;
  busy: boolean;
  feedRef: React.RefObject<HTMLDivElement>;
  supportRanking: RoomState['agents'];
  topPerformer?: RoomState['agents'][number];
  audienceMessage: string;
  audienceKind: string;
  targetAgentId: string;
  setAudienceMessage: (value: string) => void;
  setAudienceKind: (value: string) => void;
  setTargetAgentId: (value: string) => void;
  setForm: (next: RoomForm) => void;
  onCreateRoom: () => Promise<void>;
  onStartBattle: () => Promise<void>;
  onStopBattle: () => Promise<void>;
  onSendAudienceAction: (message: string, kind?: string) => Promise<void>;
};

// BattleStage 只负责“节目舞台”这一层界面：
// 包括角色席位、实时聊天、舞台摘要卡片和观众输入区。
// 数据拉取与持久化都放在外层，这样它可以保持纯 UI 组件属性。
export default function BattleStage({
  boot,
  state,
  form,
  error,
  busy,
  feedRef,
  supportRanking,
  topPerformer,
  audienceMessage,
  audienceKind,
  targetAgentId,
  setAudienceMessage,
  setAudienceKind,
  setTargetAgentId,
  setForm,
  onCreateRoom,
  onStartBattle,
  onStopBattle,
  onSendAudienceAction,
}: BattleStageProps) {
  const selectedNames = state.agents
    .filter((agent) => form.personaIds.includes(agent.id))
    .map((agent) => agent.name);

  const battleTags = (
    <Space wrap>
      <Tag bordered={false}>模式 {modeLabel(boot.modes, state.mode)}</Tag>
      <Tag bordered={false}>引擎 {state.engine}</Tag>
      <Tag bordered={false}>配置 {state.providerPath ? 'YAML 已加载' : 'YAML 未加载'}</Tag>
      <Tag color="volcano">热度 {state.heat}</Tag>
      <Tag color="gold">戏剧值 {state.dramaLevel}</Tag>
    </Space>
  );

  const stagePhase = stagePhaseLabel(state.status, state.currentRound, state.totalRounds);
  const stageProgress = state.totalRounds > 0
    ? Math.min(100, Math.max(0, Math.round((state.currentRound / state.totalRounds) * 100)))
    : 0;

  return (
    <>
      <Row gutter={[12, 12]} className="metric-row battle-metric-row">
        <Col xs={24} sm={12} xl={6}><MetricCard label="房间状态" value={statusLabel(state.status)}/></Col>
        <Col xs={24} sm={12} xl={6}><MetricCard label="当前回合" value={`${state.currentRound}/${state.totalRounds}`}/></Col>
        <Col xs={24} sm={12} xl={6}><MetricCard label="热度值" value={`${state.heat}`}/></Col>
        <Col xs={24} sm={12} xl={6}><MetricCard label="戏剧张力" value={`${state.dramaLevel}`}/></Col>
        <Col xs={24} sm={12} xl={6}><MetricCard label="支持率榜首" value={state.supportLeader || '待定'}/></Col>
      </Row>

      <div className="battle-stage-shell">
        <Sider width={320} breakpoint="lg" collapsedWidth="0" className="battle-sider">
          <Suspense fallback={<PanelFallback title="角色席位" description="正在同步选手状态。"/>}>
            <AgentRosterPanel agents={state.agents}/>
          </Suspense>
        </Sider>

        <Content className="battle-content">
          <Row gutter={[18, 18]}>
            <Col xs={24} xl={16}>
              <Card className="battle-panel battle-panel-main" bordered={false}>
                <Flex justify="space-between" align="flex-start" gap={16} wrap="wrap">
                  <div>
                    <Text className="panel-kicker">主聊天室</Text>
                    <Title level={2}>{state.topic || '请先输入本场辩题'}</Title>
                  </div>
                  {battleTags}
                </Flex>

                <Card bordered={false} className="broadcast-strip">
                  <Flex justify="space-between" align="center" wrap="wrap" gap={14}>
                    <div>
                      <Text className="panel-kicker">节目流程</Text>
                      <Title level={4}>{stagePhase}</Title>
                      <Paragraph className="broadcast-strip-copy">
                        {state.status === 'running'
                          ? `第 ${state.currentRound || 1} 回合正在播出，导播正把镜头切给 ${state.currentSpeaker || '下一位发言者'}。`
                          : state.status === 'finished'
                            ? '全场发言已经结束，舞台进入裁判复盘与赛后点评。'
                            : '房间还在准备阶段，先设置题目和上场人格。'}
                      </Paragraph>
                    </div>
                    <div className="broadcast-progress">
                      <Flex justify="space-between">
                        <Text type="secondary">播出进度</Text>
                        <Text>{stageProgress}%</Text>
                      </Flex>
                      <Progress percent={stageProgress} showInfo={false} strokeColor={{'0%': '#70a1ff', '100%': '#ff6b57'}}/>
                    </div>
                  </Flex>
                </Card>

                <Row gutter={[14, 14]} className="showcase-row">
                  <Col xs={24} md={8}>
                    <Card bordered={false} className="showcase-card showcase-card-hot">
                      <Text className="panel-kicker">火力源</Text>
                      <Title level={4}>{state.currentSpeaker || '等待开战'}</Title>
                      <Paragraph>{state.currentTarget ? `这次主要冲着 ${state.currentTarget} 去。` : state.providerReady ? '当前最可能把群聊气氛带进爆表区的麦位。' : '先把真实模型接好，这里才会开始滚动。'}</Paragraph>
                    </Card>
                  </Col>
                  <Col xs={24} md={8}>
                    <Card bordered={false} className="showcase-card">
                      <Text className="panel-kicker">默认人格</Text>
                      <Title level={4}>{selectedNames.length} 位已排上场</Title>
                      <Paragraph>{selectedNames.join(' · ') || '尚未指定'}</Paragraph>
                    </Card>
                  </Col>
                  <Col xs={24} md={8}>
                    <Card bordered={false} className="showcase-card showcase-card-cool">
                      <Text className="panel-kicker">观众带节奏</Text>
                      <Title level={4}>{state.audienceQueue} 条待处理</Title>
                      <Paragraph>{state.audienceMood}</Paragraph>
                    </Card>
                  </Col>
                </Row>

                <Card className="stage-banner" bordered={false}>
                  <Flex justify="space-between" align="center" wrap="wrap" gap={16} className="stage-banner-layout">
                    <div className="stage-banner-copy">
                      <Text className="panel-kicker">导播提示</Text>
                      <Title level={4}>{state.lastNotice}</Title>
                      <Paragraph className="stage-banner-mood">{state.audienceMood}</Paragraph>
                      <Row gutter={[12, 12]} className="stage-pressure-row">
                        <Col xs={24} md={12}>
                          <div className="stage-meter">
                            <Flex justify="space-between">
                              <Text type="secondary">全场热度</Text>
                              <Text>{state.heat}/100</Text>
                            </Flex>
                            <Progress percent={state.heat} showInfo={false} strokeColor={{'0%': '#ff8f6b', '100%': '#ff4d6d'}}/>
                          </div>
                        </Col>
                        <Col xs={24} md={12}>
                          <div className="stage-meter">
                            <Flex justify="space-between">
                              <Text type="secondary">节目张力</Text>
                              <Text>{state.dramaLevel}/100</Text>
                            </Flex>
                            <Progress percent={state.dramaLevel} showInfo={false} strokeColor={{'0%': '#ffd36b', '100%': '#ff9f43'}}/>
                          </div>
                        </Col>
                      </Row>
                    </div>
                    <Badge.Ribbon text="当前火力点" color="volcano">
                      <Card size="small" className={`speaker-card ${state.currentSpeaker ? 'speaker-card-live' : ''}`}>
                        <Text strong>{state.currentSpeaker || '等待开战'}</Text>
                        <Paragraph className="speaker-card-subline">
                          {state.currentTarget ? `正在围攻：${state.currentTarget}` : '尚未锁定攻击目标'}
                        </Paragraph>
                      </Card>
                    </Badge.Ribbon>
                  </Flex>
                </Card>

                <ChatFeed messages={state.messages} feedRef={feedRef}/>

                <Card className="composer-box" bordered={false}>
                  <Space wrap className="full-width quick-row">
                    {quickActions.map((item) => (
                      <Button key={item.label} onClick={() => void onSendAudienceAction(item.label, item.kind)} disabled={busy}>
                        {item.label}
                      </Button>
                    ))}
                  </Space>
                  <Row gutter={[14, 14]}>
                    <Col xs={24} lg={16}>
                      <TextArea
                        value={audienceMessage}
                        onChange={(event) => setAudienceMessage(event.target.value)}
                        placeholder="插一句嘴，例如：别吵了，上 benchmark。"
                        autoSize={{minRows: 4, maxRows: 6}}
                      />
                    </Col>
                    <Col xs={24} lg={8}>
                      <Space direction="vertical" size={10} className="full-width">
                        <Select<string> value={audienceKind} onChange={setAudienceKind} options={[
                          {value: '弹幕', label: '弹幕'},
                          {value: '指令', label: '指令'},
                          {value: '主持人', label: '主持人'},
                          {value: '导播', label: '导播'},
                        ]}/>
                        <Select
                          value={targetAgentId || undefined}
                          onChange={(value?: string) => setTargetAgentId(value || '')}
                          allowClear
                          placeholder="随机转火对象"
                          options={state.agents.map((agent) => ({value: agent.id, label: agent.name}))}
                        />
                        <Flex gap={10}>
                          <Button type="primary" onClick={() => void onSendAudienceAction(audienceMessage)} disabled={busy} block>
                            发送插嘴
                          </Button>
                          <Button onClick={() => void onStopBattle()} disabled={busy || state.status !== 'running'} block>
                            暂停
                          </Button>
                        </Flex>
                      </Space>
                    </Col>
                  </Row>
                </Card>
              </Card>
            </Col>

            <Col xs={24} xl={8}>
              <Suspense fallback={<PanelFallback title="导播台" description="正在装载控制台与裁判席。"/>}>
                <RightRail
                  boot={boot}
                  state={state}
                  form={form}
                  setForm={setForm}
                  busy={busy}
                  error={error}
                  topPerformer={topPerformer}
                  supportRanking={supportRanking}
                  onCreateRoom={onCreateRoom}
                  onStartBattle={onStartBattle}
                />
              </Suspense>
            </Col>
          </Row>
        </Content>
      </div>
    </>
  );
}
