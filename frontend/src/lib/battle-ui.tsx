// 文件说明：battle-ui 提供战斗页共享的展示层工具。
// 职责：集中维护指标卡片、加载占位以及模式和状态的中文展示映射。
import {Card, Statistic, Typography} from 'antd';
import type {ModeOption} from '../types';

const {Paragraph, Text, Title} = Typography;

// 这些工具函数只处理展示层细节，
// 方便把“UI 文案 / 占位态”与异步控制逻辑清楚分开。
// MetricCard 用于头部或舞台概览区的单指标展示卡片。
export function MetricCard({label, value}: {label: string; value: string}) {
  return (
    <Card bordered={false} className="metric-card">
      <Statistic title={label} value={value}/>
    </Card>
  );
}

// PanelFallback 用于 Suspense 异步组件加载中的统一占位面板。
export function PanelFallback({title, description}: {title: string; description: string}) {
  return (
    <Card className="battle-panel" bordered={false}>
      <Text className="panel-kicker">模块装载中</Text>
      <Title level={3}>{title}</Title>
      <Paragraph>{description}</Paragraph>
    </Card>
  );
}

// statusLabel 把房间状态码映射为中文显示文案。
export function statusLabel(status: string) {
  switch (status) {
    case 'running':
      return '激战中';
    case 'ready':
      return '待开场';
    case 'finished':
      return '已收官';
    case 'stopped':
      return '已暂停';
    default:
      return '未开始';
  }
}

// modeLabel 根据模式 ID 获取对应的中文标签。
export function modeLabel(modes: ModeOption[], id: string) {
  return modes.find((item) => item.id === id)?.label ?? id;
}

// formatError 把未知错误统一整理成可展示的中文字符串。
export function formatError(error: unknown) {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === 'string') {
    return error;
  }
  return '发生未知错误';
}

// stagePhaseLabel 把房间状态和回合信息映射成节目流程阶段文案。
export function stagePhaseLabel(status: string, currentRound: number, totalRounds: number) {
  if (status === 'finished') {
    return '终局复盘';
  }
  if (status === 'running' && currentRound <= 1) {
    return '开场互喷';
  }
  if (status === 'running' && currentRound < totalRounds) {
    return '中盘拉扯';
  }
  if (status === 'running') {
    return '决胜回合';
  }
  if (status === 'stopped') {
    return '临时切播';
  }
  if (status === 'ready') {
    return '候场装麦';
  }
  return '片头预热';
}
