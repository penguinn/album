package mongo

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

var (
	mongodb *mgo.Database
)

type MongoConfig struct {
	Addresses []string
	UserName  string
	Password  string
	Database  string
}

type MongoComponent struct{}

func (MongoComponent) Init(ops ...interface{}) (err error) {
	if len(ops) == 0 {
		errStr := "初始化mongo错误"
		log.Println("[mongo] error:", errStr)
		return errors.New(errStr)
	}
	cfg := ops[0].(*MongoConfig)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    cfg.Addresses,
		Username: cfg.UserName,
		Password: cfg.Password,
		Timeout:  2 * time.Second,
	})
	if err != nil {
		fmt.Println("[mongo] error:", err)
		return nil
	}
	mongodb = session.DB(cfg.Database)
	return nil
}

func Get() *mgo.Database {
	return mongodb
}
