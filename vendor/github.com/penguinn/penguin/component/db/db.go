package db

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

var (
	rander = rand.New(rand.NewSource(time.Now().Unix()))

	dbHolder = map[string]*Wrapper{}

	errDBNotFound = errors.New("DB Not Found")

	errConfig = errors.New("DBConfig Error")
)

type DBConfig map[string]struct {
	Driver string
	Source string
	Slave  map[string]struct{ Source string }
}

type DBComponent struct{}

func (DBComponent) Init(options ...interface{}) (err error) {

	if len(options) == 0 {
		return errors.New("初始化数据库错误")
	}

	c, ok := options[0].(*DBConfig)
	if !ok {
		err = errConfig
		return
	}

	for name, config := range *c {
		w := new(Wrapper)
		w.dsn, err = gorm.Open(config.Driver, config.Source)
		if err != nil {
			return
		}
		for _, s := range config.Slave {
			var slave *gorm.DB

			slave, err = gorm.Open(config.Driver, s.Source)
			if err != nil {
				return
			}
			w.slave = append(w.slave, slave)
		}

		dbHolder[name] = w
	}

	registerCallback()
	return nil
}

type Wrapper struct {
	dsn   *gorm.DB
	slave []*gorm.DB
}

func (db *Wrapper) Write() *gorm.DB {
	return db.dsn
}

func (db *Wrapper) Read() *gorm.DB {
	if len(db.slave) == 0 {
		return db.Write()
	}
	return db.slave[rander.Intn(len(db.slave))]
}

func Read(name string) (*gorm.DB, error) {
	if w, err := get(name); err == nil {
		return w.Read(), nil
	} else {
		return nil, err
	}
}

func Write(name string) (*gorm.DB, error) {
	if w, err := get(name); err == nil {
		return w.Write(), nil
	} else {
		return nil, err
	}
}

func MustRead(name string) *gorm.DB {
	return mustGet(name).Read()
}

func MustWrite(name string) *gorm.DB {
	return mustGet(name).Write()
}

func mustGet(name string) *Wrapper {
	c, ok := dbHolder[name]
	if !ok {
		panic(errDBNotFound)
	}

	return c
}

func get(name string) (*Wrapper, error) {

	c, ok := dbHolder[name]
	if !ok {
		return nil, errDBNotFound
	}

	return c, nil
}

func registerCallback() {
	gorm.DefaultCallback.Create().After("gorm:update_time_stamp").Register("my:update_time_stamp", func(scope *gorm.Scope) {
		if !scope.HasError() {
			now := time.Now().Unix()
			if ct, ok := scope.FieldByName("CreateTime"); ok {
				ct.Set(now)
			}
			if ct, ok := scope.FieldByName("UpdateTime"); ok {
				ct.Set(now)
			}
		}
	})
	gorm.DefaultCallback.Update().After("gorm:update_time_stamp").Register("my:update_time_stamp", func(scope *gorm.Scope) {
		if _, ok := scope.Get("gorm:update_column"); !ok {
			scope.SetColumn("UpdateTime", time.Now().Unix())
		}
	})
}

func NotDeletedScope(db *gorm.DB) *gorm.DB {
	return db.Where("del_flag = ?", 0)
}

func IsError(err error) error {
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
