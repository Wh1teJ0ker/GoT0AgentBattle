package main

import (
	"context"

	"GoT0AgentBattle/internal/battle"
)

// App 是暴露给前端的 Wails 桥接层。
// 它本身尽量保持轻量，只做上下文绑定和方法转发。
type App struct {
	ctx     context.Context
	service *battle.Service
}

// NewApp 创建桌面应用实例，并初始化核心业务服务。
func NewApp() *App {
	return &App{
		service: battle.NewService(),
	}
}

// startup 在 Wails 启动后触发，用于把运行时上下文注入业务层。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.service.BindContext(ctx)
}

// Bootstrap 返回首屏渲染需要的完整初始化数据。
func (a *App) Bootstrap() battle.BootstrapData {
	return a.service.Bootstrap()
}

// GetState 返回当前房间状态快照。
func (a *App) GetState() battle.RoomState {
	return a.service.State()
}

// GetSettings 返回本地持久化设置。
func (a *App) GetSettings() battle.AppSettings {
	return a.service.Settings()
}

// SaveSettings 保存设置中心的配置结果。
func (a *App) SaveSettings(input battle.AppSettings) (battle.AppSettings, error) {
	return a.service.SaveSettings(input)
}

// GenerateRandomPersona 生成一个新的随机人格模板。
func (a *App) GenerateRandomPersona() battle.PersonaConfig {
	return a.service.GenerateRandomPersona()
}

// CreateRoom 根据前端配置创建新的辩论房间。
func (a *App) CreateRoom(input battle.RoomConfigInput) battle.RoomState {
	return a.service.CreateRoom(input)
}

// StartBattle 启动当前房间的真实模型辩论流程。
func (a *App) StartBattle() (battle.RoomState, error) {
	return a.service.Start()
}

// StopBattle 手动停止当前辩论。
func (a *App) StopBattle() battle.RoomState {
	return a.service.Stop()
}

// SendAudienceAction 把观众插嘴或导播指令送入房间队列。
func (a *App) SendAudienceAction(input battle.AudienceActionInput) (battle.RoomState, error) {
	return a.service.QueueAudienceAction(input)
}
