package battle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Logger 负责把应用运行中的基础事实信息追加写入本地日志文件。
type Logger struct {
	mu   sync.Mutex
	path string
}

// NewLogger 创建一个轻量本地日志器。
// 这里故意只记录运行中的基础事实信息，不承担完整审计职责。
func NewLogger() *Logger {
	return &Logger{path: filepath.Join("data", "logs", "app.log")}
}

// Path 返回当前日志文件路径，便于外部展示或排查。
func (l *Logger) Path() string {
	return l.path
}

// Logf 以追加方式写入一行基础日志。
// 日志内容只保留对本地排障有帮助的运行痕迹，避免噪音过多。
func (l *Logger) Logf(scope string, format string, args ...any) {
	if l == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return
	}

	line := fmt.Sprintf(
		"%s [%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ToUpper(strings.TrimSpace(scope)),
		fmt.Sprintf(format, args...),
	)

	file, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	_, _ = file.WriteString(line)
}
