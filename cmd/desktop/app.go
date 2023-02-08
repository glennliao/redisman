package main

import (
	"context"
	"github.com/glennliao/apijson-go/action"
	"github.com/glennliao/apijson-go/config"
	"github.com/glennliao/apijson-go/query"
	"github.com/glennliao/redisman/server"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"net/http"
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
	var handler func(ctx context.Context, p g.Map) (g.Map, error)
	switch req.Url {
	case "/get":
		handler = Get
	case "/post":
		handler = Post
	case "/put":
		handler = Put
	case "/delete":
		handler = Delete
	}

	data, err := handler(a.ctx, req.Data.Map())
	if err != nil {
		res.Code = 500
		res.Msg = err.Error()
		g.Log().Error(a.ctx, err)
	} else {
		res.Data = data
	}
	return res
}

func Get(ctx context.Context, req g.Map) (res g.Map, err error) {
	q := query.New(ctx, req)
	q.AccessVerify = config.AccessVerify
	q.AccessCondition = config.AccessConditionFunc
	return q.Result()
}

func Head(ctx context.Context, req g.Map) (res g.Map, err error) {
	return nil, err
}

func Post(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodPost, req)
	return act.Result()
}

func Put(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodPut, req)
	return act.Result()
}

func Delete(ctx context.Context, req g.Map) (res g.Map, err error) {
	act := action.New(ctx, http.MethodDelete, req)
	return act.Result()
}
