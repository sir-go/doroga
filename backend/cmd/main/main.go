package main

import (
	"flag"
	"net/http"

	zlog "github.com/rs/zerolog/log"
)

func main() {
	// parse -c parameter
	fCfgPath := flag.String("c", "config.toml",
		"path to conf file")
	flag.Parse()

	// load the config file
	config := LoadConfig(*fCfgPath)

	// init DB
	db := new(DB)
	if err := db.Connect(config.Db); err != nil {
		panic(err)
	}

	// init Web server
	srv := &Server{dbc: db, cfg: config}

	// set teardown callbacks
	initInterrupt(db.Disconnect, srv.Shutdown)

	// run server loop
	if err := srv.Run(config.Service.Host, config.Service.Port); err != nil {
		if err != http.ErrServerClosed {
			zlog.Err(err).Msg("server running")
			panic(err)
		}
	}
}
