package ws

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/redis/go-redis/v9"
)

type Req struct {
	Id     string
	Action string
	Params *gjson.Json
	User   *User
}

type User struct {
	Id          string
	ConnectAt   *gtime.Time
	User        string
	RedisId     int
	RedisClient *redis.Client
}

type Action func(ctx context.Context, req *Req, reply func(ctx context.Context, ret any, err error))

var ActionMap = map[string]Action{}

func RegAction(name string, action Action) {
	ActionMap[name] = action
}