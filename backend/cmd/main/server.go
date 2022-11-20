package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	srv *http.Server
	dbc *DB
	cfg *Config
}

func (s *Server) Run(host string, port int) error {
	//gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(
		gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s %s%3d%s %s%-7s%s %s [%d %s]\n%s",
				param.TimeStamp.Format("06/01/02 15:04:05"),
				param.StatusCodeColor(), param.StatusCode, param.ResetColor(),
				param.MethodColor(), param.Method, param.ResetColor(),
				param.Path,
				param.BodySize,
				param.Latency,
				param.ErrorMessage,
			)
		}),
		gin.Recovery(),
		cors.Default(),
	)

	r.POST("/", s.FormProcessing)

	ds := r.Group("ds")
	// /ds
	ds.GET("/", s.GetDocumentsAll)

	ds.GET("/pub", s.GetAllDocumentsPub)

	// /ds/byid/5e5d01f67520c27b0a94f774
	ds.GET("/byid/:id", s.GetDocuments)

	arch := ds.Group("archive")
	// /ds/archive
	arch.GET("/", s.GetTarballsIndex)

	// /ds/archive/2020-03-25.tar
	arch.GET("/:filename", s.GetDocumentsTarball)

	addr := fmt.Sprintf("%s:%d", host, port)

	s.srv = &http.Server{
		Addr:           addr,
		WriteTimeout:   s.cfg.Service.Timeouts.Write.Duration,
		ReadTimeout:    s.cfg.Service.Timeouts.Read.Duration,
		IdleTimeout:    s.cfg.Service.Timeouts.Idle.Duration,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	//goland:noinspection HttpUrlsUsage
	zlog.Info().Str("running addr", "http://"+addr).Msg("")

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	zlog.Debug().Msg("shutdown web-server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
