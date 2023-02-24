package main

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/framework/handler"
	_ "github.com/glennliao/redisman/cmd/server/packed"
	"github.com/glennliao/redisman/server"
	"github.com/glennliao/redisman/server/api"
	"github.com/glennliao/redisman/server/version"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

func main() {

	version.VersionAction()

	server.Init()

	config.AccessVerify = false
	config.AccessConditionFunc = func(ctx context.Context, req config.AccessConditionReq) (g.Map, error) {
		return map[string]interface{}{}, nil
	}

	s := gins.Server()

	s.BindHandler("/ws", api.Ws)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.CORSDefault()
			r.Middleware.Next()
		})

		group.Bind(api.Http)

		handler.Bind(group)
	})

	if gres.Contains("dist/assets") {
		s.AddStaticPath("/", "dist")
		s.AddStaticPath("/assets", "dist/assets")
	}

	s.SetServerRoot(gfile.MainPkgPath())
	s.Run()
}
