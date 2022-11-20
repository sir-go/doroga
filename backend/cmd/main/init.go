package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func initInterrupt(tearDowns ...func() error) {
	zlog.Info().Msg("-- start --")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(c chan os.Signal) {
		for range c {
			for _, td := range tearDowns {
				if err := td(); err != nil {
					zlog.Err(err).Msg("shutdown")
				}
			}
			zlog.Info().Msg("-- stop --")
			os.Exit(137)
		}
	}(c)
}

func init() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})
}
