package config

import (
	"fmt"
	"github.com/luohaocode/Go-000/Week04/ent"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Option func(db *dbOptions) error

type dbOptions struct {
	timeout int64
}

var defaultDBOptions = dbOptions{
	timeout: 100,
}

func ApplyOptions(x *DbConfig) []Option {
	opts := make([]Option, 0)
	if x.Timeout != nil {
		opts = append(opts, func(db *dbOptions) error {
			db.timeout = x.Timeout.Value
			return nil
		})
	}
	return opts
}

type DBClient struct {
	addr     string
	name     string
	username string
	password string
	opts     dbOptions
}

func New(addr string, name string, username string, password string, options ...Option) (*DBClient, error) {
	e := &DBClient{
		addr:     addr,
		name:     name,
		username: username,
		password: password,
	}
	opts := defaultDBOptions
	for _, o := range options {
		err := o(&opts)
		if err != nil {
			panic(err)
		}
	}
	e.opts = opts
	return e, nil
}

func Connect(c *DBClient) (client *ent.Client, err error) {
	client, err = ent.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", c.username, c.password, c.addr, c.name))
	if err != nil {
		log.Fatal(err)
	}

	return
}
