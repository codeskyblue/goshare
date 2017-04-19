package comtool

import (
	"os"
	"os/signal"
)

func TrapSignal(f func(sig os.Signal), sigs ...os.Signal) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)
	go func() {
		for sig := range sigCh {
			f(sig)
		}
	}()
}
