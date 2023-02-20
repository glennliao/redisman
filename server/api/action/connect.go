package action

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/glennliao/redisman/server/service"
	"github.com/glennliao/redisman/server/service/model"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

func init() {
	ws.RegAction("redisConn", Conn)
}

func Conn(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {

	param := gjson.New(req.Params)
	id := param.Get("id").Int()

	connection, err := service.RedisConnection().Get(ctx, id)
	if err != nil {
		reply(ctx, nil, err)
		return
	}

	// todo 检查id权限
	// 获取conn meta

	if req.User.RedisId != id {

	} else {

	}

	if req.User.RedisClient != nil {
		req.User.RedisClient.Close()
	}
	if req.User.LocalListener != nil {
		req.User.LocalListener.Close()
	}

	client, err := getRedisClientByConfig(ctx, req, connection)

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

func getRedisClientByConfig(ctx context.Context, req *ws.Req, connection *model.RedisConnectionGet) (client *redis.Client, err error) {

	redisOptions := &redis.Options{
		Addr:        net.JoinHostPort(connection.Host, gconv.String(connection.Port)),
		Username:    connection.Username,
		Password:    connection.Password,
		DB:          gconv.Int(connection.Db),
		DialTimeout: time.Second * 5,
		PoolSize:    1,
	}

	// tls
	tlsOption := connection.Options.Tls
	if tlsOption.Enable {
		cfg, err := getTLSConfig(tlsOption)
		if err != nil {
			return nil, gerror.Wrap(err, "parse cert and key")
		}
		redisOptions.TLSConfig = cfg
	}

	// ssh
	sshOption := connection.Options.Ssh
	if sshOption.Enable {

		redisOptions.ReadTimeout = -2
		redisOptions.WriteTimeout = -2

		sshConfig, err := getSSHConfig(sshOption)
		if err != nil {
			return nil, err
		}

		// ssh + tls
		if tlsOption.Enable {
			port, err := sshTunnel(ctx, req, sshConfig, sshOption, net.JoinHostPort(connection.Host, gconv.String(connection.Port)))
			redisOptions.Addr = net.JoinHostPort(connection.Host, gconv.String(port))
			if err != nil {
				return nil, gerror.Wrap(err, "ssh sshTunnel")
			}
		} else {
			sshClient, err := ssh.Dial("tcp", net.JoinHostPort(sshOption.Host, gconv.String(sshOption.Port)), sshConfig)
			if err != nil {
				return nil, gerror.Wrap(err, "ssh conn")
			}

			redisOptions.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return sshClient.Dial(network, addr)
			}
		}
	}

	client = redis.NewClient(redisOptions)

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, gerror.Wrap(err, "redis connect")
	}

	return
}

func getSSHConfig(sshOption model.SSH) (sshConfig *ssh.ClientConfig, err error) {
	var authMethods []ssh.AuthMethod
	if sshOption.AuthType == "password" {
		authMethods = append(authMethods, ssh.Password(sshOption.Password))
	} else {
		signer, err := signerFromPem([]byte(sshOption.PrivateKey), []byte(sshOption.Passphrase))
		if err != nil {
			return nil, gerror.Wrap(err, "ssh private")
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	return &ssh.ClientConfig{
		User:            sshOption.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}, nil
}

func getTLSConfig(tlsOption model.TLS) (tlsConfig *tls.Config, err error) {
	tlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
		//ServerName:   "you domain",
		InsecureSkipVerify: true,
	}

	certificate, err := tls.X509KeyPair([]byte(tlsOption.Cert), []byte(tlsOption.Key))
	if err != nil {
		return nil, gerror.Wrap(err, "parse cert and key")
	}
	tlsConfig.Certificates = []tls.Certificate{certificate}

	if tlsOption.Ca != "" {
		ca := x509.NewCertPool()
		ok := ca.AppendCertsFromPEM([]byte(tlsOption.Ca))
		if !ok {
			return nil, gerror.New("ca parse fail")
		}
		tlsConfig.RootCAs = ca
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

func sshTunnel(ctx context.Context, req *ws.Req, config *ssh.ClientConfig, sshOption model.SSH, redisAddr string) (port int, err error) {

	localListener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, gerror.Wrap(err, "net.Listen failed")
	}

	go g.TryCatch(ctx, func(ctx context.Context) {
		for {

			localConn, err := localListener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					g.Log().Error(ctx, "localListener.Accept failed", err)
				}
				break
			}

			sshClientConn, err := ssh.Dial("tcp", net.JoinHostPort(sshOption.Host, sshOption.Port), config)
			if err != nil {
				g.Log().Error(nil, "ssh.Dial failed", err)
				return
			}

			sshConn, err := sshClientConn.Dial("tcp", redisAddr)

			if err != nil {
				g.Log().Error(ctx, "sshClientConn failed", err, redisAddr)
				continue
			}

			// 将localConn.Reader复制到sshConn.Writer
			go func() {
				_, err = io.Copy(sshConn, localConn)
				if err != nil {
					g.Log().Error(ctx, "io.Copy failed", err)
					return
				}
			}()

			// 将sshConn.Reader复制到localConn.Writer
			go func() {
				_, err = io.Copy(localConn, sshConn)
				if err != nil {
					g.Log().Error(ctx, "io.Copy failed", err)
					return
				}
			}()

		}
	}, func(ctx context.Context, exception error) {
		g.Log().Error(ctx, exception)
	})

	req.User.LocalListener = localListener

	return localListener.Addr().(*net.TCPAddr).Port, nil
}
