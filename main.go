package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// main 是桌面端应用入口。
// 这里主要负责把前端静态资源、窗口配置和 Wails 生命周期绑定起来。
func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:         "GoT0AgentBattle",
		Width:         1520,
		Height:        960,
		MinWidth:      1280,
		MinHeight:     820,
		Frameless:     false,
		DisableResize: false,
		StartHidden:   false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 7, G: 10, B: 17, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
