export namespace battle {
	
	export class AgentState {
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
	    status: string;
	    model: string;
	    modelProfileId: string;
	    lastLine: string;
	    lastTarget: string;
	    currentTurn: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AgentState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.role = source["role"];
	        this.avatar = source["avatar"];
	        this.color = source["color"];
	        this.style = source["style"];
	        this.persona = source["persona"];
	        this.tagline = source["tagline"];
	        this.aggressive = source["aggressive"];
	        this.toxicity = source["toxicity"];
	        this.enabled = source["enabled"];
	        this.anger = source["anger"];
	        this.support = source["support"];
	        this.tokenUsage = source["tokenUsage"];
	        this.status = source["status"];
	        this.model = source["model"];
	        this.modelProfileId = source["modelProfileId"];
	        this.lastLine = source["lastLine"];
	        this.lastTarget = source["lastTarget"];
	        this.currentTurn = source["currentTurn"];
	    }
	}
	export class PersonaConfig {
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
	
	    static createFrom(source: any = {}) {
	        return new PersonaConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.role = source["role"];
	        this.avatar = source["avatar"];
	        this.color = source["color"];
	        this.style = source["style"];
	        this.persona = source["persona"];
	        this.tagline = source["tagline"];
	        this.aggressive = source["aggressive"];
	        this.toxicity = source["toxicity"];
	        this.enabled = source["enabled"];
	        this.defaultModelProfileId = source["defaultModelProfileId"];
	        this.systemPrompt = source["systemPrompt"];
	    }
	}
	export class RoomConfigInput {
	    topic: string;
	    mode: string;
	    agentCount: number;
	    model: string;
	    rounds: number;
	    providerPath: string;
	    preferRealLLM: boolean;
	    personaIds: string[];
	    judgeModelProfileId: string;
	
	    static createFrom(source: any = {}) {
	        return new RoomConfigInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.topic = source["topic"];
	        this.mode = source["mode"];
	        this.agentCount = source["agentCount"];
	        this.model = source["model"];
	        this.rounds = source["rounds"];
	        this.providerPath = source["providerPath"];
	        this.preferRealLLM = source["preferRealLLM"];
	        this.personaIds = source["personaIds"];
	        this.judgeModelProfileId = source["judgeModelProfileId"];
	    }
	}
	export class AppSettings {
	    defaultConfig: RoomConfigInput;
	    personas: PersonaConfig[];
	    judgeModelProfileId: string;
	    preferredThemeFlavor: string;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultConfig = this.convertValues(source["defaultConfig"], RoomConfigInput);
	        this.personas = this.convertValues(source["personas"], PersonaConfig);
	        this.judgeModelProfileId = source["judgeModelProfileId"];
	        this.preferredThemeFlavor = source["preferredThemeFlavor"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AudienceActionInput {
	    message: string;
	    targetAgentId: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new AudienceActionInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.targetAgentId = source["targetAgentId"];
	        this.kind = source["kind"];
	    }
	}
	export class JudgeSummary {
	    round: number;
	    winnerAgentId: string;
	    winnerName: string;
	    winnerReason: string;
	    winnerSide: string;
	    effectiveArguments: string[];
	    invalidTrashTalk: string[];
	    showmanshipScore: number;
	    entertainmentRating: string;
	
	    static createFrom(source: any = {}) {
	        return new JudgeSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.round = source["round"];
	        this.winnerAgentId = source["winnerAgentId"];
	        this.winnerName = source["winnerName"];
	        this.winnerReason = source["winnerReason"];
	        this.winnerSide = source["winnerSide"];
	        this.effectiveArguments = source["effectiveArguments"];
	        this.invalidTrashTalk = source["invalidTrashTalk"];
	        this.showmanshipScore = source["showmanshipScore"];
	        this.entertainmentRating = source["entertainmentRating"];
	    }
	}
	export class DebateMessage {
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
	    performanceTag: string;
	
	    static createFrom(source: any = {}) {
	        return new DebateMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.agentId = source["agentId"];
	        this.agentName = source["agentName"];
	        this.agentAvatar = source["agentAvatar"];
	        this.color = source["color"];
	        this.content = source["content"];
	        this.replyTo = source["replyTo"];
	        this.mentions = source["mentions"];
	        this.timestamp = source["timestamp"];
	        this.round = source["round"];
	        this.tone = source["tone"];
	        this.isJudge = source["isJudge"];
	        this.isAudience = source["isAudience"];
	        this.heatDelta = source["heatDelta"];
	        this.supportDelta = source["supportDelta"];
	        this.performanceTag = source["performanceTag"];
	    }
	}
	export class RoomState {
	    roomId: string;
	    topic: string;
	    mode: string;
	    model: string;
	    status: string;
	    currentRound: number;
	    totalRounds: number;
	    heat: number;
	    supportLeader: string;
	    currentSpeaker: string;
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
	    judgeSummary?: JudgeSummary;
	
	    static createFrom(source: any = {}) {
	        return new RoomState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.roomId = source["roomId"];
	        this.topic = source["topic"];
	        this.mode = source["mode"];
	        this.model = source["model"];
	        this.status = source["status"];
	        this.currentRound = source["currentRound"];
	        this.totalRounds = source["totalRounds"];
	        this.heat = source["heat"];
	        this.supportLeader = source["supportLeader"];
	        this.currentSpeaker = source["currentSpeaker"];
	        this.lastNotice = source["lastNotice"];
	        this.audienceQueue = source["audienceQueue"];
	        this.startedAt = source["startedAt"];
	        this.finishedAt = source["finishedAt"];
	        this.engine = source["engine"];
	        this.providerPath = source["providerPath"];
	        this.providerReady = source["providerReady"];
	        this.providerNotice = source["providerNotice"];
	        this.judgeModelProfileId = source["judgeModelProfileId"];
	        this.personaIds = source["personaIds"];
	        this.agents = this.convertValues(source["agents"], AgentState);
	        this.messages = this.convertValues(source["messages"], DebateMessage);
	        this.judgeSummary = this.convertValues(source["judgeSummary"], JudgeSummary);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ModelProfileSummary {
	    id: string;
	    label: string;
	    model: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelProfileSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.model = source["model"];
	        this.description = source["description"];
	    }
	}
	export class ModelOption {
	    id: string;
	    label: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.description = source["description"];
	    }
	}
	export class ModeOption {
	    id: string;
	    label: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new ModeOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.description = source["description"];
	    }
	}
	export class BootstrapData {
	    modes: ModeOption[];
	    models: ModelOption[];
	    modelProfiles: ModelProfileSummary[];
	    defaultConfig: RoomConfigInput;
	    settings: AppSettings;
	    settingsPath: string;
	    state: RoomState;
	    realProviderReady: boolean;
	    realProviderPath: string;
	    realProviderNotice: string;
	
	    static createFrom(source: any = {}) {
	        return new BootstrapData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modes = this.convertValues(source["modes"], ModeOption);
	        this.models = this.convertValues(source["models"], ModelOption);
	        this.modelProfiles = this.convertValues(source["modelProfiles"], ModelProfileSummary);
	        this.defaultConfig = this.convertValues(source["defaultConfig"], RoomConfigInput);
	        this.settings = this.convertValues(source["settings"], AppSettings);
	        this.settingsPath = source["settingsPath"];
	        this.state = this.convertValues(source["state"], RoomState);
	        this.realProviderReady = source["realProviderReady"];
	        this.realProviderPath = source["realProviderPath"];
	        this.realProviderNotice = source["realProviderNotice"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	

}

