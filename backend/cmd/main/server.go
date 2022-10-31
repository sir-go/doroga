package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Run() {
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

	r.POST("/", SRV.FormProcessing)

	// /ds
	ds := r.Group("ds")
	ds.GET("/", SRV.GetDSall)

	ds.GET("/pub", SRV.GetDSallPub)
	// /ds/byid/5e5d01f67520c27b0a94f774
	ds.GET("/byid/:id", SRV.GetDS)

	// /ds/archive
	arch := ds.Group("archive")
	arch.GET("/", SRV.GetDSTarIndex)
	// /ds/archive/2020-03-25.tar
	arch.GET("/:filename", SRV.GetDSTar)

	addr := fmt.Sprintf("%s:%d", CFG.Service.Host, CFG.Service.Port)
	LOG.Println("run web-server on http://" + addr)

	s.srv = &http.Server{
		Addr:           addr,
		WriteTimeout:   CFG.Service.Timeouts.Write.Duration,
		ReadTimeout:    CFG.Service.Timeouts.Read.Duration,
		IdleTimeout:    CFG.Service.Timeouts.Idle.Duration,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		if err != nil {
			LOG.Panic(err)
		}
	}
}

func (s *Server) Shutdown() {
	LOG.Println("shutdown web-server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		LOG.Panic(err)
	}
}
