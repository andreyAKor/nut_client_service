package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/andreyAKor/nut_parser/internal/app"
	"github.com/andreyAKor/nut_parser/internal/configs"
	clientsNut "github.com/andreyAKor/nut_parser/internal/http/clients/nut"
	"github.com/andreyAKor/nut_parser/internal/http/server"
	"github.com/andreyAKor/nut_parser/internal/logging"
	metricsNut "github.com/andreyAKor/nut_parser/internal/metrics/nut"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "nut_parser",
	Short: "NUT parser service application",
	Long:  "The NUT parser service is the most simplified service for storing resized images from another web-resources and storing this resized images in own file lru-cache.",
	RunE:  run,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")
	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:funlen
func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	c := new(configs.Config)
	if err := c.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}

	// Init clients
	nutClient, err := clientsNut.New(c.Clients.NUT.Host, c.Clients.NUT.Port, c.Clients.NUT.Username, c.Clients.NUT.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize NUT client")
	}

	// Init http-server
	srv, err := server.New(c.HTTP.Host, c.HTTP.Port, c.HTTP.BodyLimit, nutClient)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-server")
	}

	// Init metrics
	nutMetrics, err := metricsNut.New(c.Metrics.NUT.Interval, nutClient)

	// Init and run app
	a, err := app.New(srv, nutMetrics)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize app")
	}
	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("app runnign fail")
	}

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	<-interruptCh

	log.Info().Msg("Stopping...")

	if err := srv.Close(); err != nil {
		log.Fatal().Err(err).Msg("http-server closing fail")
	}
	if err := a.Close(); err != nil {
		log.Fatal().Err(err).Msg("app closing fail")
	}

	log.Info().Msg("Stopped")

	if err := l.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
