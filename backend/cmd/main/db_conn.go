package main

import (
	"github.com/globalsign/mgo"
	zlog "github.com/rs/zerolog/log"
)

type DB struct {
	collectionName string
	session        *mgo.Session
}

func (db *DB) Connect(cfg CfgDb) (err error) {
	db.session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{cfg.Host},
		Timeout:  cfg.Timeout.Duration,
		Database: cfg.DbName,
		Username: cfg.User,
		Password: cfg.Password})
	if err != nil {
		zlog.Err(err).Msg("db connect")
	}
	if err == nil {
		db.session.SetMode(mgo.Monotonic, true)
	}
	db.collectionName = cfg.Collection
	return
}

func (db *DB) Disconnect() error {
	zlog.Debug().Msg("shutdown db connection")
	db.session.Close()
	return nil
}
