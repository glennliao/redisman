package model

import "github.com/gogf/gf/v2/encoding/gjson"

type SSH struct {
	Enable     bool
	Host       string
	Port       string
	Username   string
	Password   string
	AuthType   string // password / private_key
	PrivateKey string
	Passphrase string
}

type TLS struct {
	Enable bool
	Cert   string
	Key    string
	Ca     string
}

type Options struct {
	Ssh SSH
	Tls TLS
}

type RedisConnectionGet struct {
	Id       int
	Title    string
	Host     string
	Port     string
	Db       string
	Username string
	Password string
	DbAlias  *gjson.Json
	Options  Options
}
