package main

import (
	"context"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/framework/handler"
	"github.com/glennliao/apijson-go/model"
	"github.com/glennliao/redisman/server"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

// App struct
type App struct {
	ctx  context.Context
	user ws.User
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	server.Init()

	id := guid.S()
	a.user = ws.User{
		Id:        id,
		ConnectAt: gtime.Now(),
	}

	config.AccessVerify = false
	config.AccessConditionFunc = func(ctx context.Context, req config.AccessConditionReq) (g.Map, error) {
		return map[string]interface{}{}, nil
	}
}

func (a *App) Action(msg string) string {
	var req ws.Req
	err := gconv.Scan(msg, &req)
	if err != nil {
		panic(err)
	}

	var retMap = map[string]any{
		"id":   req.Id,
		"code": 200,
	}

	req.User = &a.user

	ws.Handler(a.ctx, &req, func(ctx context.Context, ret any, err error) {
		retMap["data"] = ret
		if err != nil {
			retMap["code"] = 500
			retMap["msg"] = err.Error()
			g.Log().Error(ctx, err)
		}
	})

	return gjson.New(retMap).MustToJsonString()
}

type HttpReq struct {
	Method string
	Url    string
	Data   *gjson.Json
}

type HttpRes struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func (a *App) Http(req HttpReq) HttpRes {
	res := HttpRes{Code: 200, Data: req.Data}
	var httpHandler func(ctx context.Context, p model.Map) (model.Map, error)
	switch req.Url {
	case "/get":
		httpHandler = handler.Get
	case "/post":
		httpHandler = handler.Post
	case "/put":
		httpHandler = handler.Put
	case "/delete":
		httpHandler = handler.Delete
	}

	data, err := httpHandler(a.ctx, req.Data.Map())
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		g.Log().Error(a.ctx, err)
	} else {
		res.Data = data
	}
	return res
}
