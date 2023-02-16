package api

import (
	"context"
	_ "github.com/glennliao/redisman/server/api/action"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
	"sync"
)

var users = sync.Map{}

func Ws(r *ghttp.Request) {

	var ctx = r.Context()
	conn, err := r.WebSocket()
	if err != nil {
		glog.Error(ctx, err)
		r.Exit()
	}

	id := guid.S()
	user := ws.User{
		Id:        id,
		ConnectAt: gtime.Now(),
	}
	users.Store(id, &user)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				g.Log().Error(ctx, err)
				break
			}
		}

		var req ws.Req
		err = gconv.Scan(msg, &req)
		if err != nil {
			g.Log().Error(ctx, err)
			break
		}

		req.User = &user

		ws.Handler(ctx, &req, func(ctx context.Context, ret any, err error) {
			var retMap = map[string]any{
				"id":   req.Id,
				"data": ret,
				"code": 200,
			}

			if err != nil {
				retMap["code"] = 500
				retMap["msg"] = err.Error()
				g.Log().Error(ctx, err)
			}

			if err = conn.WriteMessage(msgType, gjson.New(retMap).MustToJson()); err != nil {
				return
			}
		})
	}

	if user.RedisClient != nil {
		user.RedisClient.Close()
	}
	users.Delete(id)

}
