package aclapi

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)	
	}
}

func exec() error {
	
	/* config must load here in exec() if needed to */

	utils.InitLogger(true)

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
