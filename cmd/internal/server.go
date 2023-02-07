package internal

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/config/db"
	"github.com/glennliao/apijson-go/framework"
	"github.com/glennliao/apijson-go/framework/handler"
	"github.com/glennliao/redisman/api"
	_ "github.com/glennliao/redisman/cmd/server/packed"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/samber/lo"
)

func Server() {

	gres.Dump()

	initConfig()

	config.AccessVerify = false
	config.AccessConditionFunc = func(ctx context.Context, req config.AccessConditionReq) (g.Map, error) {
		return map[string]interface{}{}, nil
	}

	initDb()

	db.Reload()

	s := g.Server()

	s.BindHandler("/ws", api.Ws)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.CORSDefault()
			r.Middleware.Next()
		})

		group.Bind(api.Http)

		handler.Bind(group)
	})

	s.AddStaticPath("/", "dist")
	s.AddStaticPath("/assets", "dist/assets")

	s.SetServerRoot(gfile.MainPkgPath())
	s.Run()
}

func initConfig() {
	if !gfile.Exists("./config.yaml") {
		gfile.PutContents("./config.yaml", "server:\n  address: \":16379\"\n  dumpRouterMap: false\n\nlogger:\n  default:\n    level: all\n\ndatabase:\n  logger:\n    level: \"all\"\n    stdout: true\n  default:\n    debug: true\n#    link: \"mysql:root:fLkYC4t&9wLX6&cU@tcp(192.168.32.71:3306)/redis_test?charset=utf8mb4&parseTime=True&loc=Local\"\n    link: \"sqlite::@file(./db.sqlite3)?charset=utf8mb4&parseTime=True&loc=Local\"\n\n")
	}
}

// 先这么写着吧
func initDb() {
	redis_connection_mysql := "CREATE TABLE `redis_connection` (\n                                    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,\n                                    `title` varchar(255) DEFAULT NULL,\n                                    `host` varchar(255) DEFAULT NULL,\n                                    `port` int(5) DEFAULT NULL,\n                                    `db` int(5) DEFAULT 0 NULL,\n                                    `username` varchar(128) DEFAULT NULL,\n                                    `password` varchar(128) DEFAULT NULL,\n                                    `options` json DEFAULT NULL,\n                                    `db_alias` json DEFAULT NULL,\n                                    `tags` json DEFAULT NULL,\n                                    `created_at` datetime DEFAULT NULL,\n                                    `created_by` varchar(32) DEFAULT NULL,\n                                    `updated_at` datetime DEFAULT NULL,\n                                    `updated_by` varchar(32) DEFAULT NULL,\n                                    `deleted_at` datetime DEFAULT NULL,\n                                    PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"
	redis_connection_sqlite := "CREATE TABLE `redis_connection` (\n                                    `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,\n                                    `title` varchar(255) DEFAULT NULL,\n                                    `host` varchar(255) DEFAULT NULL,\n                                    `port` varchar(5) DEFAULT NULL,\n                                    `db` varchar(5) DEFAULT '0' NULL,\n                                    `username` varchar(128) DEFAULT NULL,\n                                    `password` varchar(128) DEFAULT NULL,\n                                    `options` varchar(5) DEFAULT NULL,\n                                    `db_alias` text DEFAULT NULL,\n                                    `tags` varchar(512) DEFAULT NULL,\n                                    `created_at` datetime DEFAULT NULL,\n                                    `created_by` varchar(32) DEFAULT NULL,\n                                    `updated_at` datetime DEFAULT NULL,\n                                    `updated_by` varchar(32) DEFAULT NULL,\n                                    `deleted_at` datetime DEFAULT NULL\n                \n) ;"

	ctx := gctx.New()

	tables, err := g.DB().Tables(ctx)
	if err != nil {
		panic(err)
	}
	if !lo.Contains(tables, "redis_connection") {
		dbType := g.DB().GetConfig().Type

		redis_connection_table := ""

		switch dbType {
		case "sqlite":
			redis_connection_table = redis_connection_sqlite
			sql_access := "CREATE TABLE `_access` (\n  `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,\n  `debug` tinyint(4) NOT NULL DEFAULT '0' ,\n  `name` varchar(50) NOT NULL ,\n  `alias` varchar(20) DEFAULT NULL ,\n  `get` varchar(100) NOT NULL  ,\n  `head` varchar(100) NOT NULL  ,\n  `gets` varchar(100) NOT NULL  ,\n  `heads` varchar(100) NOT NULL  ,\n  `post` varchar(100) NOT NULL  ,\n  `put` varchar(100) NOT NULL  ,\n  `delete` varchar(100) NOT NULL  ,\n  `created_at` datetime NOT NULL ,\n  `detail` varchar(1000) DEFAULT NULL,\n  `row_key` varchar(32) DEFAULT NULL ,\n  `fields_get` text DEFAULT NULL ,\n  `row_key_gen` varchar(255) DEFAULT NULL,\n  `executor` varchar(32) DEFAULT NULL\n);"
			sql_request := "CREATE TABLE `_request` (\n  `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,\n  `debug` tinyint(4) NOT NULL DEFAULT '0' ,\n  `version` tinyint(4) NOT NULL DEFAULT '1' ,\n  `method` varchar(10) DEFAULT 'GETS' ,\n  `tag` varchar(20) NOT NULL ,\n  `structure` text NOT NULL ,\n  `detail` varchar(10000) DEFAULT NULL ,\n  `created_at` datetime,\n  `exec_queue` varchar(255) DEFAULT NULL ,\n  `executor` text DEFAULT NULL \n);"
			sql_function := "CREATE TABLE `_function` (\n  `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,\n  `debug` tinyint(4) NOT NULL DEFAULT '0' ,\n  `userId` bigint(20) NOT NULL ,\n  `name` varchar(50) NOT NULL ,\n  `arguments` varchar(100) DEFAULT NULL ,\n  `demo` text ,\n  `detail` varchar(1000) NOT NULL ,\n  `type` varchar(50) NOT NULL DEFAULT 'Object' ,\n  `version` tinyint(4) NOT NULL DEFAULT '0' ,\n  `tag` varchar(20) DEFAULT NULL ,\n  `methods` varchar(50) DEFAULT NULL ,\n  `created_at` datetime NOT NULL ,\n  `back` varchar(45) DEFAULT NULL \n);"

			g.DB().Exec(ctx, sql_access)
			g.DB().Exec(ctx, sql_request)
			g.DB().Exec(ctx, sql_function)

		case "mysql":
			redis_connection_table = redis_connection_mysql
			framework.Init()

		}
		g.DB().Exec(ctx, redis_connection_table)

		g.DB().Insert(ctx, "_access", g.Map{
			"debug":   0,
			"name":    "redis_connection",
			"alias":   "RedisConnection",
			"get":     "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"head":    "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"gets":    "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"heads":   "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"post":    "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"put":     "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"delete":  "[\"UNKNOWN\", \"LOGIN\", \"CONTACT\", \"CIRCLE\", \"OWNER\", \"ADMIN\"]",
			"row_key": "id",
		})

		g.DB().Insert(ctx, "_request", g.Map{
			"debug":     0,
			"version":   "1",
			"tag":       "RedisConnection",
			"method":    "POST",
			"structure": "{}",
		})

		g.DB().Insert(ctx, "_request", g.Map{
			"debug":     0,
			"version":   "1",
			"tag":       "RedisConnection",
			"method":    "PUT",
			"structure": "{}",
		})

		g.DB().Insert(ctx, "_request", g.Map{
			"debug":     0,
			"version":   "1",
			"tag":       "RedisConnection",
			"method":    "DELETE",
			"structure": "{}",
		})

	}

}
