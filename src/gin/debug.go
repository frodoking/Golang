package gin

import (
	"log"
)

func init() {
	log.SetFlags(0)
}

func IsDebugging() bool {
	return ginMode == debugCode
}

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
	if IsDebugging() {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers[nuHandlers-1])
		debugPrint("%-5s %-25s --> %s (%d handlers) \n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[GIN-debug] "+format, values...)
	}
}
func debugPrintWARNING() {
	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
 	- using env:	export GIN_MODE=release
 	- using code:	gin.SetMode(gin.ReleaseMode)
	`)
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}
