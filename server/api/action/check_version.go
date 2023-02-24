package action

import (
	"context"
	"github.com/glennliao/redisman/server/api/ws"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	ws.RegAction("check-version", checkVersion)
}

type ReleaseInfo struct {
	Version string
	Remark  map[string]string
}

const CountryCN = "CN"
const LatestReleasePrefix = "https://github.com/glennliao/redisman/releases/latest/download/"
const GithubProxy = "https://ghproxy.com/"

const ReleaseInfoJson = "release_info.json"

func checkVersion(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {
	country, err := getCountry(ctx)
	if err != nil {
		reply(ctx, nil, err)
		return
	}

	urlPrefix := LatestReleasePrefix
	if country == CountryCN { // proxy for china
		urlPrefix = GithubProxy + urlPrefix
	}

	info, err := checkLatestVersion(ctx, urlPrefix)
	if err != nil {
		reply(ctx, nil, err)
		return
	}

	reply(ctx, g.Map{"latest": info}, nil)

}

func checkLatestVersion(ctx context.Context, urlPrefix string) (info *ReleaseInfo, err error) {
	resp, err := gclient.New().Get(ctx, urlPrefix+ReleaseInfoJson)
	if err != nil {
		return nil, err
	}

	err = gconv.Scan(resp.ReadAllString(), &info)
	if err != nil {
		return nil, err
	}

	return info, err
}

func getCountry(ctx context.Context) (country string, err error) {
	res, err := gclient.New().Get(ctx, "https://ipinfo.io/json")
	if err != nil {
		return "", err
	}
	type IpInfo struct {
		Country string
	}
	var info IpInfo
	err = gconv.Scan(res.ReadAllString(), &info)
	return info.Country, err
}
