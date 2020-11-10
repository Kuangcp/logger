package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type brush func(string) string

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

//鉴于终端的通常使用习惯，一般白色和黑色字体是不可行的,所以30,37不可用，
var colors = []brush{
	newBrush("1;41"), // Emergency          红色底
	newBrush("1;35"), // Alert              紫色
	newBrush("1;34"), // Critical           蓝色
	newBrush("1;31"), // Error              红色
	newBrush("1;33"), // Warn               黄色
	newBrush("1;36"), // Informational      天蓝色
	newBrush("1;32"), // Debug              绿色
	newBrush("1;32"), // Trace              绿色
}

type ConsoleLogger struct {
	sync.Mutex
	Level    string `json:"level"`
	Colorful bool   `json:"color"`
	logLevel int
}

func (c *ConsoleLogger) Init(debug bool, jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	if jsonConfig != "{}" && debug {
		fmt.Fprintf(os.Stdout, "ConsoleLogger Init:%s\n", jsonConfig)
	}

	err := json.Unmarshal([]byte(jsonConfig), c)
	if runtime.GOOS == "windows" {
		c.Colorful = false
	}

	if l, ok := LevelMap[c.Level]; ok {
		c.logLevel = l
		return nil
	}

	return err
}

func (c *ConsoleLogger) LogWrite(when time.Time, msgText interface{}, level int) error {
	if level > c.logLevel {
		return nil
	}
	msg, ok := msgText.(string)
	if !ok {
		return nil
	}
	if c.Colorful {
		msg = colors[level](msg)
	}
	c.printlnConsole(when, msg)
	return nil
}

func (c *ConsoleLogger) Destroy() {

}

func (c *ConsoleLogger) printlnConsole(when time.Time, msg string) {
	c.Lock()
	defer c.Unlock()
	os.Stdout.Write(append([]byte(msg), '\n'))
}

func init() {
	Register(AdapterConsole, &ConsoleLogger{
		logLevel: LevelDebug,
		Colorful: runtime.GOOS != "windows",
	})
}
