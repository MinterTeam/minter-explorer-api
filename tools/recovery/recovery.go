package recovery

import (
	log "github.com/sirupsen/logrus"
	"runtime/debug"
)

func SafeGo(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.WithField("stacktrace", string(debug.Stack())).Error(err)
		}
	}()

	f()
}
