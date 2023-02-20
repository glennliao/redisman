package server

import (
	"github.com/glennliao/apijson-go/framework"
	"github.com/glennliao/apijson-go/framework/database"
	_ "github.com/glennliao/redisman/server/api/action"
	"github.com/glennliao/table-sync/tablesync"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/samber/lo"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Init() {
	InitConfig()
	initDb()
	framework.Init()
}

func InitConfig() {

	homeDir, _ := os.UserHomeDir()

	configPath := filepath.Join(homeDir, "/.config/redisman/config.toml")

	configContent := `
[server]
    address = ":16379"
    dumpRouterMap = false
[logger]
    path = "./log/"
    level = "all"
[database]
    [database.logger]
        Level = "all"
        Stdout = true
    [database.default]
        Debug = true
        link = "sqlite::@file(./db.sqlite3)"

	`

	configContent = strings.ReplaceAll(configContent, "./db.sqlite3", path.Join(homeDir, "/.config/redisman/db.sqlite3"))
	configContent = strings.ReplaceAll(configContent, "\\", "/")

	fileAdapter := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	fileAdapter.AddPath(gfile.Dir(configPath))
	if _, err := fileAdapter.GetFilePath(); err != nil {
		gfile.PutContents(configPath, configContent)
	}

}

func initDb() {

	ctx := gctx.New()

	_db := g.DB()

	tables, err := _db.Tables(ctx)
	if err != nil {
		panic(err)
	}

	syncer := tablesync.Syncer{
		Tables: []tablesync.Table{
			RedisConnection{},
			database.Access{},
			database.Request{},
		},
	}

	err = syncer.Sync(ctx, _db)
	if err != nil {
		panic(err)
	}

	if !lo.Contains(tables, "redis_connection") {

		_db.Insert(ctx, "_access", map[string]any{
			"debug":   0,
			"name":    "redis_connection",
			"alias":   "RedisConnection",
			"get":     "LOGIN,OWNER,ADMIN",
			"head":    "LOGIN,OWNER,ADMIN",
			"gets":    "LOGIN,OWNER,ADMIN",
			"heads":   "LOGIN,OWNER,ADMIN",
			"post":    "LOGIN,OWNER,ADMIN",
			"put":     "OWNER,ADMIN",
			"delete":  "OWNER,ADMIN",
			"row_key": "id",
		})

		_db.Insert(ctx, "_request", []any{
			map[string]any{
				"debug":     0,
				"version":   "1",
				"tag":       "RedisConnection",
				"method":    "POST",
				"structure": "{}",
			},
			map[string]any{
				"debug":     0,
				"version":   "1",
				"tag":       "RedisConnection",
				"method":    "PUT",
				"structure": "{}",
			},
			map[string]any{
				"debug":     0,
				"version":   "1",
				"tag":       "RedisConnection",
				"method":    "DELETE",
				"structure": "{}",
			},
		})

	}

}
