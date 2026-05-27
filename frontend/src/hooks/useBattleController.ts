// 文件说明：useBattleController 负责收敛前端主流程状态。
// 职责：统一管理 Wails 调用、事件订阅、本地视图状态和页面派生数据。
import {useEffect, useMemo, useRef, useState} from 'react';
import {EventsOff, EventsOn} from '../../wailsjs/runtime/runtime';
import {
  Bootstrap,
  CreateRoom,
  GenerateRandomPersona,
  GetState,
  SaveSettings,
  SendAudienceAction,
  StartBattle,
  StopBattle,
} from '../../wailsjs/go/main/App';
import {battle as battleModels} from '../../wailsjs/go/models';
import type {AppSettings, BootstrapData, EventPayload, PersonaConfig, RoomForm, RoomState} from '../types';
import {formatError} from '../lib/battle-ui';

const EVENT_NAME = 'got0:broadcast';

function castModel<T>(value: unknown) {
  return value as T;
}

// useBattleController 统一收敛 Wails 调用、页面本地状态和派生视图数据。
// 这样 App.tsx 就能尽量保持声明式布局，review 时不必在主入口里追异步细节。
export function useBattleController() {
  const [boot, setBoot] = useState<BootstrapData | null>(null);
  const [state, setState] = useState<RoomState | null>(null);
  const [form, setForm] = useState<RoomForm | null>(null);
  const [settings, setSettings] = useState<AppSettings | null>(null);
  const [view, setView] = useState<'battle' | 'settings'>('battle');
  const [audienceMessage, setAudienceMessage] = useState('');
  const [targetAgentId, setTargetAgentId] = useState('');
  const [audienceKind, setAudienceKind] = useState('弹幕');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const [saveBusy, setSaveBusy] = useState(false);
  const feedRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function init() {
      try {
        const [bootstrap, currentState] = await Promise.all([Bootstrap(), GetState()]);
        if (cancelled) {
          return;
        }
        const typedBootstrap = castModel<BootstrapData>(bootstrap);
        setBoot(typedBootstrap);
        setState(castModel<RoomState>(currentState));
        setSettings(typedBootstrap.settings);
        setForm(typedBootstrap.defaultConfig);
      } catch (err) {
        if (!cancelled) {
          setError(formatError(err));
        }
      }
    }

    void init();

    EventsOn(EVENT_NAME, (payload: EventPayload) => {
      setState(castModel<RoomState>(payload.state));
      if (payload.notice) {
        setError('');
      }
    });

    return () => {
      cancelled = true;
      EventsOff(EVENT_NAME);
    };
  }, []);

  useEffect(() => {
    if (!feedRef.current) {
      return;
    }
    // 聊天区默认始终跟随到最新消息。
    // 这个页面更像实时舞台，而不是需要用户中途翻历史消息的 IM 聊天窗口。
    feedRef.current.scrollTop = feedRef.current.scrollHeight;
  }, [state?.messages.length]);

  useEffect(() => {
    if (!settings || form) {
      return;
    }
    setForm(settings.defaultConfig);
  }, [settings, form]);

  const supportRanking = useMemo(() => {
    const agents = state?.agents ?? [];
    return [...agents].sort((a, b) => b.support - a.support);
  }, [state?.agents]);

  const topPerformer = supportRanking[0];

  async function refreshBootstrapState() {
    const bootstrap = await Bootstrap();
    const typedBootstrap = castModel<BootstrapData>(bootstrap);
    setBoot(typedBootstrap);
    setState(typedBootstrap.state);
    return typedBootstrap;
  }

  async function createRoom() {
    if (!form) {
      return;
    }
    if (!form.topic.trim()) {
      setError('请先输入本场辩论主题。');
      return;
    }
    setBusy(true);
    setError('');
    try {
      const next = await CreateRoom({
        topic: form.topic,
        mode: form.mode,
        agentCount: Number(form.agentCount),
        model: form.model,
        rounds: Number(form.rounds),
        providerPath: form.providerPath ?? '',
        preferRealLLM: Boolean(form.preferRealLLM),
        personaIds: form.personaIds,
        judgeModelProfileId: form.judgeModelProfileId,
      });
      setState(castModel<RoomState>(next));
      setTargetAgentId('');
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  }

  async function startBattle() {
    setBusy(true);
    setError('');
    try {
      const next = await StartBattle();
      setState(castModel<RoomState>(next));
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  }

  async function stopBattle() {
    setBusy(true);
    setError('');
    try {
      const next = await StopBattle();
      setState(castModel<RoomState>(next));
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  }

  async function persistSettings() {
    if (!settings) {
      return;
    }
    setSaveBusy(true);
    setError('');
    try {
      const saved = await SaveSettings(new battleModels.AppSettings(settings));
      const typedSaved = castModel<AppSettings>(saved);
      setSettings(typedSaved);
      setForm(typedSaved.defaultConfig);
      const typedBootstrap = await refreshBootstrapState();
      setBoot(typedBootstrap);
    } catch (err) {
      setError(formatError(err));
    } finally {
      setSaveBusy(false);
    }
  }

  async function generateRandomPersona() {
    if (!settings) {
      return;
    }
    setBusy(true);
    setError('');
    try {
      const persona = castModel<PersonaConfig>(await GenerateRandomPersona());
      setSettings({
        ...settings,
        personas: settings.personas.concat(persona),
      });
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  }

  async function sendAudienceAction(message: string, kind = audienceKind) {
    const trimmed = message.trim();
    if (!trimmed) {
      return;
    }
    setBusy(true);
    setError('');
    try {
      const next = await SendAudienceAction({
        message: trimmed,
        targetAgentId,
        kind,
      });
      setState(castModel<RoomState>(next));
      setAudienceMessage('');
    } catch (err) {
      setError(formatError(err));
    } finally {
      setBusy(false);
    }
  }

  return {
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
  };
}
