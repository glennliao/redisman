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

	var connection *model.RedisConnectionGet
	err := gjson.New(req.Params).Scan(&connection)
	if err != nil {
		reply(ctx, nil, gerror.Cause(err))
		return
	}

	client, err := getRedisClientByConfig(ctx, connection)

	if err != nil {
		reply(ctx, nil, err)
		return
	}

	defer func() {
		client.Close()
	}()

	reply(ctx, g.Map{}, nil)

}
