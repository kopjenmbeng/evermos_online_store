package infrastructure

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx/reflectx"
	_ "fmt"
	_ "time"
	"strings"
	log "github.com/sirupsen/logrus"
)

type IPostgresSql interface {
	Write() *sqlx.DB
	Read() *sqlx.DB
}

type PostgresSql struct {
	log   *log.Logger
	read  string
	write string
}

func NewPostgresSql(read string,write string,log *log.Logger)IPostgresSql{
	return &PostgresSql{read: read,write: write,log: log}
}

func(pg *PostgresSql)Write()*sqlx.DB{
	db := sqlx.MustOpen("postgres", pg.write)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	// ensure the database are reachable
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	pg.log.Println("write db connected successfully !")
	return db
}

func(pg *PostgresSql)Read()*sqlx.DB{
	db := sqlx.MustOpen("postgres", pg.read)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	// ensure the database are reachable
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	pg.log.Println("read db connected successfully !")
	return db
}