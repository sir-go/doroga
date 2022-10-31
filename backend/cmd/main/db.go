package main

import (
	"github.com/globalsign/mgo"
)

type DB struct {
	session *mgo.Session
}

func (db *DB) Connect() (err error) {
	db.session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{CFG.Db.Host},
		Timeout:  CFG.Db.Timeout.Duration,
		Database: CFG.Db.DbName,
		Username: CFG.Db.User,
		Password: CFG.Db.Password})
	if err == nil {
		db.session.SetMode(mgo.Monotonic, true)
	}
	return err
}

func (db *DB) Disconnect() {
	db.session.Close()
}
