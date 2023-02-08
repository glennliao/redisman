package server

import (
	"github.com/glennliao/table-sync/tablesync"
	"time"
)

type RedisConnection struct {
	tablesync.TableMeta
	Id        uint32 `ddl:"primaryKey"`
	Title     string
	Host      string
	Port      string `ddl:"size:5"`
	Db        string `ddl:"size:2;default:0"`
	Username  string `ddl:"size:128"`
	Password  string `ddl:"size:128"`
	Options   string `ddl:"type:json"`
	DbAlias   string `ddl:"type:json"`
	Tags      string `ddl:"type:json"`
	CreatedAt *time.Time
	CreatedBy string `ddl:"size:32"`
	UpdatedAt *time.Time
	UpdatedBy string `ddl:"size:32"`
	DeletedAt *time.Time
}
