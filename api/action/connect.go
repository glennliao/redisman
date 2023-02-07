package action

import (
	"context"
	"fmt"
	"github.com/glennliao/redisman/api/ws"
	"github.com/glennliao/redisman/service"
	"github.com/glennliao/redisman/service/model"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func init() {
	ws.RegAction("redisConn", Conn)
}

func Conn(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {

	param := gjson.New(req.Params)
	id := param.Get("id").Int()

	g.Dump(req.Params)

	connection, err := service.RedisConnection().Get(ctx, id)
	if err != nil {
		reply(ctx, nil, err)
		return
	}

	g.Dump(connection)

	// todo 检查id权限
	// 获取conn meta

	if req.User.RedisId != id {

	} else {

	}

	if req.User.RedisClient != nil {
		req.User.RedisClient.Close()
	}

	client, err := getRedisClientByConfig(ctx, connection)

	if err != nil {
		reply(ctx, nil, err)
		return
	}

	req.User.RedisClient = client
	req.User.RedisId = id

	reply(ctx, g.Map{
		"title":   connection.Title,
		"dbAlias": connection.DbAlias,
	}, nil)
}

func getRedisClientByConfig(ctx context.Context, connection *model.RedisConnectionGet) (client *redis.Client, err error) {
	var dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	if connection.Options.Ssh.Enable {
		sshOption := connection.Options.Ssh

		var authMethods []ssh.AuthMethod
		if sshOption.AuthType == "password" {
			authMethods = append(authMethods, ssh.Password(sshOption.Password))
		} else {
			signer, err := signerFromPem([]byte(sshOption.PrivateKey), []byte(sshOption.Passphrase))
			if err != nil {
				return nil, gerror.Cause(err)
			}
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}

		sshConfig := &ssh.ClientConfig{
			User:            sshOption.Username,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         15 * time.Second,
		}

		sshClient, err := ssh.Dial("tcp", net.JoinHostPort(sshOption.Host, gconv.String(sshOption.Port)), sshConfig)
		if err != nil {
			return nil, gerror.Cause(err)
		}

		dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		}
	}

	client = redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(connection.Host, gconv.String(connection.Port)),
		Username:     connection.Username,
		Password:     connection.Password,
		DB:           gconv.Int(connection.Db),
		DialTimeout:  time.Second * 5,
		Dialer:       dialer,
		ReadTimeout:  -2,
		WriteTimeout: -2,
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, gerror.Cause(err)
	}

	return
}

func signerFromPem(pemBytes []byte, passphrase []byte) (signer ssh.Signer, err error) {

	if len(passphrase) > 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, passphrase)
	} else {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	}

	if err != nil {
		return nil, fmt.Errorf("parsing plain private key failed %v", err)
	}

	return signer, nil
}
