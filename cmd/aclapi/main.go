package aclapi

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/PythonHacker24/linux-acl-management-aclapi"
	"go.uber.org/zap"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)	
	}
}

func exec() error {
	
	/* config must load here in exec() if needed to */

	utils.InitLogger(true)

	/* zap.L() can be used all over the code for global level logging */
	zap.L().Info("Logger Initiated ...")

	/* preparing graceful shutdown */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupt
		cancel()
	}()

	return run(ctx)
}

func run(ctx context.Context) error {
		

}
