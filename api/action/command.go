package action

import (
	"context"
	"github.com/glennliao/redisman/api/ws"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
)

func init() {
	ws.RegAction("redisCom", Command)
}

func Command(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {

	var rdb = req.User.RedisClient
	pipe := rdb.Pipeline()

	var commands [][]any
	err := req.Params.Scan(&commands)
	if err != nil {
		reply(ctx, nil, err)
		return
	}

	size := len(commands)
	retCmd := make([]*redis.Cmd, size)
	rets := make([]interface{}, size)

	for i := 0; i < size; i++ {
		retCmd[i] = pipe.Do(ctx, commands[i]...)
	}

	if _, err = pipe.Exec(ctx); nil != err {
		g.Log().Warning(ctx, commands)
		g.Log().Error(ctx, err)
		reply(ctx, rets, err)
		return
	}

	for i := 0; i < len(retCmd); i++ {
		ret, err := retCmd[i].Result()

		if nil != err {
			panic(err)
		}

		if v, ok := ret.(map[any]any); ok {
			rets[i] = gconv.Map(v)
		} else {
			rets[i] = ret
		}
	}

	reply(ctx, rets, nil)

}
