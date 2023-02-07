package api

import (
	"context"
	_ "github.com/glennliao/redisman/api/action"
	"github.com/glennliao/redisman/api/ws"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
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
			//if e, ok := err.(net.Error); ok {
			//	if e.Error() == "close 1001 (going away)" {
			//		fmt.Println("close")
			//	} else {
			//		fmt.Println(e.Error())
			//	}
			//	break
			//}
			//panic(err)
			break
		}

		var req ws.Req
		err = gconv.Scan(msg, &req)
		if err != nil {
			panic(err)
			break
		}

		req.User = &user

		handler(ctx, &req, func(ctx context.Context, ret any, err error) {
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

func handler(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {
	for k, action := range ws.ActionMap {
		if k == req.Action {
			action(ctx, req, reply)
			break
		}
	}
}
