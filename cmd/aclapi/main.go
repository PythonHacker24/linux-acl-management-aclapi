package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-aclapi/config"
	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver"
	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/utils"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)
	}
}

func exec() error {

	/* config must load here in exec() if needed to */

	/* setting up cobra for cli interactions */
	var (
		configPath string
		rootCmd    = &cobra.Command{
			Use:   "aclapi <command> <subcommand>",
			Short: "API Daemon for linux acl management",
			Example: heredoc.Doc(`
				$ aclapi --config /path/to/config.yaml
			`),
			Run: func(cmd *cobra.Command, args []string) {
				if configPath != "" {
					fmt.Printf("Using config file: %s\n\n", configPath)
				} else {
					fmt.Printf("No config file provided.\n\n")
				}
			},
		}
	)

	/* adding --config argument */
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file")

	/* Execute the command */
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("arguments error: %s", err.Error())
		os.Exit(1)
	}

	/*
		load config file
		if there is an error in loading the config file, then it will exit with code 1
	*/
	if err := config.LoadConfig(configPath); err != nil {
		fmt.Printf("Configuration Error in %s: %s",
			configPath,
			err.Error(),
		)
		/* since the configuration is invalid, don't proceed */
		os.Exit(1)
	}

	/*
		true for production, false for development mode
		logger is only for gRPC server and core components (after this step)
	*/
	utils.InitLogger(!config.APIDConfig.DConfig.DebugMode)

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

	grpcServer, err := grpcserver.InitServer()
	if err != nil {
		zap.L().Error("Failed to initialize gRPC server",
			zap.Error(err),
		)
	}

	/* creating the gRPC listener */
	listener, err := grpcServer.Start()
	if err != nil {
		zap.L().Error("Failed to start gRPC server",
			zap.Error(err),
		)
		return err
	}

	/* starting the gRPC listener */
	go func() {
		zap.L().Info("gRPC server is listening",
			zap.String("Address", config.APIDConfig.Server.Host),
			zap.Int("Port", config.APIDConfig.Server.GrpcPort),
		)
		if err := grpcServer.GRPC.Serve(listener); err != nil {
			zap.L().Error("gRPC server error",
				zap.Error(err),
			)
		}
	}()

	<-ctx.Done()

	zap.L().Info("Shutting down gRPC server")

	/* attempting to gracefully shutdown gRPC server */
	gracefulStop := make(chan struct{})
	go func() {
		grpcServer.GRPC.GracefulStop()
		close(gracefulStop)
	}()

	select {
	case <-gracefulStop:
		zap.L().Info("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		/* gRPC server graceful shutdown timeout reached, need to forcefully shutdown it down */
		zap.L().Info("Graceful shutdown timeout reached. Forcing gRPC server shutdown")
		grpcServer.GRPC.Stop()
	}

	return nil
}
