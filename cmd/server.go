/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package cmd

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"xyz/cmd/rest"
	"xyz/pkg/otel"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start rest api server",
	Long:  `Initialize rest api server and run it`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		fiberApp := rest.New(ctx)
		defer func() {
			fiberApp.Shutdown()
			otel.Shutdown()
		}()

		signKill := make(chan os.Signal, 1)
		signal.Notify(signKill, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-signKill

		log.Info("\n\n=========================================\n\n")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
