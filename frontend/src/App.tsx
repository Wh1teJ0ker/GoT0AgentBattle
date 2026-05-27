// 文件说明：App 是前端正式界面的组合壳层入口。
// 职责：装配主题、顶层布局、主视图切换，并把控制器状态分发给战斗页和设置页。
import {lazy, Suspense} from 'react';
import {
  Alert,
  App as AntApp,
  Card,
  ConfigProvider,
  Flex,
  Layout,
  Segmented,
  Typography,
} from 'antd';
import './App.css';
import {PanelFallback} from './lib/battle-ui';
import {useBattleController} from './hooks/useBattleController';
import BattleStage from './components/BattleStage';

const {Header, Content} = Layout;
const {Title, Paragraph, Text} = Typography;

const SettingsPanel = lazy(() => import('./components/SettingsPanel'));

const theme = {
  token: {
    colorPrimary: '#ff6b57',
    colorInfo: '#4cc9f0',
    colorSuccess: '#2ed573',
    colorWarning: '#ff9f43',
    colorError: '#ff4d6d',
    borderRadius: 18,
    fontFamily: '"Avenir Next", "Segoe UI", sans-serif',
    colorBgBase: '#070a11',
    colorTextBase: '#f6f7fb',
  },
  components: {
    Layout: {
      bodyBg: 'transparent',
      headerBg: 'transparent',
      siderBg: 'transparent',
    },
    Card: {
      colorBgContainer: 'rgba(11, 14, 24, 0.9)',
      colorBorderSecondary: 'rgba(255,255,255,0.08)',
    },
  },
} as const;

// App 只承担组合壳层职责：
// 主题、顶层布局和主视图切换留在这里，
// 具体状态逻辑与舞台渲染拆到独立模块中。
function App() {
  const {
    boot,
    state,
    form,
    settings,
    view,
    audienceMessage,
    targetAgentId,
    audienceKind,
    error,
    busy,
    saveBusy,
    feedRef,
    supportRanking,
    topPerformer,
    setForm,
    setSettings,
    setView,
    setAudienceMessage,
    setTargetAgentId,
    setAudienceKind,
    createRoom,
    startBattle,
    stopBattle,
    persistSettings,
    generateRandomPersona,
    sendAudienceAction,
  } = useBattleController();

  if (!boot || !state || !form || !settings) {
    return (
      <div className="ant-shell">
        <Card className="loading-card">
          <Text className="loading-kicker">GoT0AgentBattle</Text>
          <Title level={1} className="loading-title">导播台启动中</Title>
          <Paragraph>正在给所有嘴替发麦。</Paragraph>
        </Card>
      </div>
    );
  }

  return (
    <ConfigProvider theme={theme}>
      <AntApp>
        <div className={`ant-shell theme-${settings.preferredThemeFlavor || 'broadcast-chaos'}`}>
          <div className="bg-grid"/>
          <div className="bg-radial bg-radial-left"/>
          <div className="bg-radial bg-radial-right"/>
          <div className="bg-spotlight bg-spotlight-left"/>
          <div className="bg-spotlight bg-spotlight-right"/>

          <Layout className="battle-layout">
            <Header className="battle-header">
              <Flex justify="space-between" align="flex-start" gap={20} wrap="wrap">
                <div>
                  <Text className="loading-kicker">AI 多角色辩论模拟器</Text>
                  <Title className="battle-title">GoT0AgentBattle</Title>
                  <Paragraph className="header-subtitle">让人格、模型和情绪值一起上桌，不讲武德地互喷。</Paragraph>
                </div>
                <Flex vertical gap={12} align="flex-end" className="header-actions">
                  <Segmented
                    value={view}
                    onChange={(value) => setView(value as 'battle' | 'settings')}
                    options={[
                      {label: '主舞台', value: 'battle'},
                      {label: '设置中心', value: 'settings'},
                    ]}
                  />
                  {view === 'settings' ? (
                    <Card bordered={false} className="metric-card settings-blurb">
                      <Text>当前共配置 {settings.personas.length} 个可编辑人格。模型能力与裁判档案由本地 YAML 配置文件统一控制。</Text>
                    </Card>
                  ) : null}
                </Flex>
              </Flex>
            </Header>

            {view === 'settings' ? (
              <Content className="battle-content">
                <Suspense fallback={<PanelFallback title="设置中心" description="正在装载人格库和总控台。"/>}>
                  <SettingsPanel
                    boot={boot}
                    settings={settings}
                    setSettings={setSettings}
                    busy={busy}
                    saveBusy={saveBusy}
                    settingsPath={boot.settingsPath}
                    onSave={persistSettings}
                    onGenerateRandomPersona={generateRandomPersona}
                  />
                </Suspense>
              </Content>
            ) : (
              <BattleStage
                boot={boot}
                state={state}
                form={form}
                error={error}
                busy={busy}
                feedRef={feedRef}
                supportRanking={supportRanking}
                topPerformer={topPerformer}
                audienceMessage={audienceMessage}
                audienceKind={audienceKind}
                targetAgentId={targetAgentId}
                setAudienceMessage={setAudienceMessage}
                setAudienceKind={setAudienceKind}
                setTargetAgentId={setTargetAgentId}
                setForm={setForm}
                onCreateRoom={createRoom}
                onStartBattle={startBattle}
                onStopBattle={stopBattle}
                onSendAudienceAction={sendAudienceAction}
              />
            )}
          </Layout>

          {error ? (
            <div className="global-error">
              <Alert type="error" showIcon message={error}/>
            </div>
          ) : null}
        </div>
      </AntApp>
    </ConfigProvider>
  );
}

export default App;
