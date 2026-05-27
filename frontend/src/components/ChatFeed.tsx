// 文件说明：ChatFeed 负责渲染实时滚动的聊天消息流。
// 职责：统一管理消息卡片的动画、语气视觉样式、裁判态与观众态展示规则。
import {Avatar, Card, Flex, Space, Tag, Typography} from 'antd';
import {memo} from 'react';
import {motion} from 'framer-motion';
import type {DebateMessage} from '../types';

const {Paragraph, Text} = Typography;

type ChatFeedProps = {
  messages: DebateMessage[];
  feedRef: React.RefObject<HTMLDivElement>;
};

// ChatFeed 单独拆出，是为了避免高频消息追加时拖着整页一起重渲染。
// 同时也把所有“某类消息为什么长这样”的视觉规则集中到一个地方，方便 review。
function ChatFeed({messages, feedRef}: ChatFeedProps) {
  return (
    <div className="chat-feed" ref={feedRef}>
      {messages.map((message, index) => (
        <motion.div
          key={message.id}
          initial={entryMotion(message)}
          animate={{opacity: 1, y: 0, x: 0, scale: 1, rotate: 0, filter: 'blur(0px)'}}
          transition={{
            duration: message.isJudge ? 0.46 : 0.28,
            ease: message.isJudge ? [0.22, 1, 0.36, 1] : 'easeOut',
            delay: Math.min(index * 0.012, 0.12),
          }}
          whileHover={{y: -2, scale: 1.004}}
          className={`chat-entry ${message.isJudge ? 'chat-entry-judge' : ''} ${message.isAudience ? 'chat-entry-audience' : ''}`}
        >
          <Card
            size="small"
            className={`chat-card ${toneClassName(message)}`}
            bordered={false}
          >
            <div className="chat-card-glow"/>
            <div className="chat-card-noise"/>
            <div className="chat-card-heatbar">
              <span style={{width: `${messageHeatWidth(message)}%`, background: messageAccent(message)}}/>
            </div>
            <Flex gap={12} align="flex-start">
              <Avatar className="chat-avatar" style={{background: message.color}}>
                {message.agentAvatar}
              </Avatar>
              <div className="chat-body">
                <Space wrap size={[8, 8]} className="chat-meta">
                  <Text strong>{message.agentName}</Text>
                  <Text type="secondary">第 {message.round || 0} 回合</Text>
                  {message.replyTo ? <Text type="secondary">@{message.replyTo}</Text> : null}
                  {message.cue ? <Tag bordered={false} className="chat-cue-tag">{message.cue}</Tag> : null}
                  <Tag bordered={false} color={messageTagColor(message)}>{message.performanceTag}</Tag>
                  {message.mentions?.map((mention) => (
                    <Tag key={`${message.id}-${mention}`} bordered={false} className="chat-mention-tag">
                      @{mention}
                    </Tag>
                  ))}
                </Space>
                <Space wrap size={[8, 8]} className="chat-submeta">
                  <Text type="secondary">冲击 {message.impact || 0}</Text>
                  <Text type="secondary">热度 +{message.heatDelta || 0}</Text>
                  <Text type="secondary">支持 +{message.supportDelta || 0}</Text>
                </Space>
                <Paragraph className="chat-content">{message.content}</Paragraph>
              </div>
            </Flex>
          </Card>
        </motion.div>
      ))}
    </div>
  );
}

// toneClassName 把后端返回的语气标签映射为前端展示分组。
// 这里不追求语义绝对精确，而是追求在高频更新下仍然有稳定的视觉分层。
function toneClassName(message: DebateMessage) {
  if (message.isJudge) {
    return 'chat-card-judge';
  }
  if (message.isAudience) {
    return 'chat-card-audience';
  }
  if (message.tone.includes('火力') || message.tone.includes('高能')) {
    return 'chat-card-hype';
  }
  if (message.tone.includes('严谨')) {
    return 'chat-card-precision';
  }
  return 'chat-card-standard';
}

// 下面这些颜色和入场辅助函数只依赖消息元数据本身，
// 这样动效表现是可预测的，review 时也更容易检查。
function messageTagColor(message: DebateMessage) {
  if (message.isJudge) {
    return 'gold';
  }
  if (message.isAudience) {
    return 'default';
  }
  if (message.cue.includes('精准') || message.cue.includes('围攻')) {
    return 'magenta';
  }
  if (message.tone.includes('火力') || message.tone.includes('高能')) {
    return 'volcano';
  }
  if (message.tone.includes('严谨')) {
    return 'blue';
  }
  return 'processing';
}

function messageAccent(message: DebateMessage) {
  if (message.isJudge) {
    return 'linear-gradient(90deg, #ffd36b, #ff9f43)';
  }
  if (message.isAudience) {
    return 'linear-gradient(90deg, #b2bec3, #dfe6e9)';
  }
  if (message.cue.includes('精准') || message.cue.includes('围攻')) {
    return 'linear-gradient(90deg, #ff78c4, #ff4d6d)';
  }
  if (message.tone.includes('火力') || message.tone.includes('高能')) {
    return 'linear-gradient(90deg, #ff6b57, #ff4d6d)';
  }
  if (message.tone.includes('严谨')) {
    return 'linear-gradient(90deg, #70a1ff, #4cc9f0)';
  }
  return `linear-gradient(90deg, ${message.color}, rgba(255,255,255,0.6))`;
}

function messageHeatWidth(message: DebateMessage) {
  const base = 24 + Math.min(52, Math.abs(message.heatDelta || 0) * 6 + Math.abs(message.supportDelta || 0) * 3);
  if (message.isJudge) {
    return 100;
  }
  if ((message.impact || 0) >= 85) {
    return 96;
  }
  if (message.isAudience) {
    return Math.max(base, 40);
  }
  return base;
}

function entryMotion(message: DebateMessage) {
  if (message.isJudge) {
    return {opacity: 0, y: 28, scale: 0.96, rotate: -0.4, filter: 'blur(8px)'};
  }
  if (message.isAudience) {
    return {opacity: 0, x: 24, y: 6, scale: 0.98, rotate: 0.4, filter: 'blur(6px)'};
  }
  return {opacity: 0, y: 16, x: -8, scale: 0.985, rotate: -0.15, filter: 'blur(4px)'};
}

export default memo(ChatFeed);
