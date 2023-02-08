package service

import (
	"context"
	"github.com/glennliao/redisman/server/service/model"
	"github.com/gogf/gf/v2/frame/g"
)

type RedisConnectionService struct {
}

var redisConnection = RedisConnectionService{}

func RedisConnection() *RedisConnectionService {
	return &redisConnection
}

func (s *RedisConnectionService) Get(ctx context.Context, id int) (row *model.RedisConnectionGet, err error) {
	err = g.DB().Model("redis_connection").Where("id", id).Scan(&row)
	return
}
