// 文件说明：AgentRosterPanel 负责渲染左侧角色席位列表。
// 职责：集中展示每个 Agent 的头像、人格标签、状态、愤怒值、支持率与 Token 估算。
import type {CSSProperties} from 'react';
import {Avatar, Card, Flex, Progress, Space, Tag, Typography} from 'antd';
import {motion} from 'framer-motion';
import type {AgentState} from '../types';

const {Paragraph, Text, Title} = Typography;

type AgentRosterPanelProps = {
  agents: AgentState[];
};

function AgentRosterPanel({agents}: AgentRosterPanelProps) {
  return (
    <Card className="battle-panel" bordered={false}>
      <Text className="panel-kicker">选手池</Text>
      <Title level={3}>角色席位</Title>
      <Space direction="vertical" size={12} className="full-width">
        {agents.map((agent, index) => (
          <motion.div
            key={agent.id}
            initial={{opacity: 0, x: -16}}
            animate={{opacity: 1, x: 0}}
            transition={{delay: index * 0.04}}
          >
            <Card
              size="small"
              className={`agent-panel ${agent.currentTurn ? 'agent-panel-live' : ''}`}
              style={{'--agent-color': agent.color} as CSSProperties}
              bordered={false}
            >
              <Flex align="center" gap={12}>
                <Avatar className="agent-avatar" style={{background: agent.color}}>
                  {agent.avatar}
                </Avatar>
                <div>
                  <Flex align="center" gap={8}>
                    <Text strong>{agent.name}</Text>
                    <Tag bordered={false}>{agent.persona}</Tag>
                  </Flex>
                  <Text type="secondary">{agent.role}</Text>
                </div>
              </Flex>
              <Paragraph className="agent-tagline">{agent.tagline}</Paragraph>
              <MiniProgress label="愤怒值" value={agent.anger} strokeColor={agent.color}/>
              <MiniProgress label="节奏值" value={agent.momentum} strokeColor="#ff8f6b"/>
              <MiniProgress label="支持率" value={agent.support} strokeColor="#feca57"/>
              <MiniProgress
                label="Token"
                value={Math.min(100, Math.round(agent.tokenUsage / 8))}
                strokeColor="#70a1ff"
              />
              <Flex justify="space-between" className="agent-meta">
                <Text type="secondary">{agent.status}</Text>
                <Text type="secondary">{agent.roastCount} 次开火</Text>
              </Flex>
              <Text type="secondary" className="agent-model-line">{agent.model}</Text>
            </Card>
          </motion.div>
        ))}
      </Space>
    </Card>
  );
}

function MiniProgress({
  label,
  value,
  strokeColor,
}: {
  label: string;
  value: number;
  strokeColor: string;
}) {
  return (
    <div className="mini-progress">
      <Flex justify="space-between">
        <Text type="secondary">{label}</Text>
        <Text>{value}</Text>
      </Flex>
      <Progress percent={value} size="small" showInfo={false} strokeColor={strokeColor}/>
    </div>
  );
}

export default AgentRosterPanel
