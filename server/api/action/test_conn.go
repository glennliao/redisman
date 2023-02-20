package action

import (
	"context"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/glennliao/redisman/server/service/model"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	ws.RegAction("redisConnTest", RedisTest)
}

func RedisTest(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {

	tempReq := &ws.Req{
		Params: req.Params,
		User:   &ws.User{},
	}

	var connection *model.RedisConnectionGet
	err := gjson.New(tempReq.Params).Scan(&connection)
	if err != nil {
		reply(ctx, nil, gerror.Wrap(err, "req scan"))
		return
	}

	client, err := getRedisClientByConfig(ctx, tempReq, connection)

	if err != nil {
		reply(ctx, nil, err)
		return
	}

	defer func() {
		client.Close()
		if tempReq.User.LocalListener != nil {
			tempReq.User.LocalListener.Close()
		}
	}()

	reply(ctx, g.Map{}, nil)

}
