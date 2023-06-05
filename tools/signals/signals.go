package signals

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func CatchStopSignal(ctx context.Context) error {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-signalChannel:
		logrus.Printf("received signal: %s", sig)
		return fmt.Errorf("received signal: %s", sig)
	case <-ctx.Done():
		return fmt.Errorf("catchStopSignal: context canceled")
	}
}
