package initial

import "github.com/jhonnli/logs"

func InitLog() {
	logs.Async()
	logs.Async(1e4)
	logs.SetLogger(logs.AdapterConsole, `{"level":7}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}
