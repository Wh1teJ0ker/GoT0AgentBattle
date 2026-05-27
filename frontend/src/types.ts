export type ModeOption = {
  id: string;
  label: string;
  description: string;
};

export type ModelOption = {
  id: string;
  label: string;
  description: string;
};

export type ModelProfileSummary = {
  id: string;
  label: string;
  model: string;
  description: string;
};

export type PersonaConfig = {
  id: string;
  name: string;
  role: string;
  avatar: string;
  color: string;
  style: string;
  persona: string;
  tagline: string;
  aggressive: number;
  toxicity: number;
  enabled: boolean;
  defaultModelProfileId: string;
  systemPrompt: string;
};

export type AgentState = {
  id: string;
  name: string;
  role: string;
  avatar: string;
  color: string;
  style: string;
  persona: string;
  tagline: string;
  aggressive: number;
  toxicity: number;
  enabled: boolean;
  anger: number;
  support: number;
  tokenUsage: number;
  momentum: number;
  roastCount: number;
  status: string;
  model: string;
  modelProfileId: string;
  lastLine: string;
  lastTarget: string;
  currentTurn: boolean;
};

export type DebateMessage = {
  id: string;
  agentId: string;
  agentName: string;
  agentAvatar: string;
  color: string;
  content: string;
  replyTo: string;
  mentions: string[];
  timestamp: number;
  round: number;
  tone: string;
  isJudge: boolean;
  isAudience: boolean;
  heatDelta: number;
  supportDelta: number;
  impact: number;
  cue: string;
  performanceTag: string;
};

export type JudgeSummary = {
  round: number;
  winnerAgentId: string;
  winnerName: string;
  winnerReason: string;
  winnerSide: string;
  effectiveArguments: string[];
  invalidTrashTalk: string[];
  showmanshipScore: number;
  entertainmentRating: string;
};

export type RoomState = {
  roomId: string;
  topic: string;
  mode: string;
  model: string;
  status: string;
  currentRound: number;
  totalRounds: number;
  heat: number;
  dramaLevel: number;
  audienceMood: string;
  supportLeader: string;
  currentSpeaker: string;
  currentTarget: string;
  lastNotice: string;
  audienceQueue: number;
  startedAt: number;
  finishedAt: number;
  engine: string;
  providerPath: string;
  providerReady: boolean;
  providerNotice: string;
  judgeModelProfileId: string;
  personaIds: string[];
  agents: AgentState[];
  messages: DebateMessage[];
  judgeSummary?: JudgeSummary | null;
};

export type AppSettings = {
  defaultConfig: RoomForm;
  personas: PersonaConfig[];
  judgeModelProfileId: string;
  preferredThemeFlavor: string;
};

export type BootstrapData = {
  modes: ModeOption[];
  models: ModelOption[];
  modelProfiles: ModelProfileSummary[];
  defaultConfig: RoomForm;
  settings: AppSettings;
  settingsPath: string;
  state: RoomState;
  realProviderReady: boolean;
  realProviderPath: string;
  realProviderNotice: string;
};

export type EventPayload = {
  type: string;
  state: RoomState;
  message?: DebateMessage;
  summary?: JudgeSummary;
  notice?: string;
  timestamp: number;
};

export type RoomForm = {
  topic: string;
  mode: string;
  agentCount: number;
  model: string;
  rounds: number;
  providerPath?: string;
  preferRealLLM?: boolean;
  personaIds: string[];
  judgeModelProfileId: string;
};
