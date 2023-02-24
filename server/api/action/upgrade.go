package action

//
//func init() {
//	ws.RegAction("upgrade", upgrade)
//}
//
//type ReleaseInfo struct {
//	Version string
//	Remark  map[string]string
//}
//
//const CountryCN = "CN"
//const LatestReleasePrefix = "https://github.com/glennliao/redisman/releases/latest/download/"
//const GithubProxy = "https://ghproxy.com/"
//
//const ReleaseInfoJson = "release_info.json"
//
//func upgrade(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {
//	country, err := getCountry(ctx)
//	if err != nil {
//		reply(ctx, nil, err)
//		return
//	}
//
//	urlPrefix := LatestReleasePrefix
//	if country == CountryCN { // proxy for china
//		urlPrefix = GithubProxy + urlPrefix
//	}
//
//	info, err := checkLatestVersion(ctx, urlPrefix)
//	if err != nil {
//		reply(ctx, nil, err)
//		return
//	}
//
//	const version = "0.3.0"
//	if info.Version > version {
//		reply(ctx, g.Map{"canUpgrade": false, "version": info.Version, "remark": info.Remark}, nil)
//		return
//	}
//
//	err = downloadNewVersion(ctx, urlPrefix)
//	if err != nil {
//		reply(ctx, nil, err)
//		return
//	}
//
//	reply(ctx, g.Map{"hadDownloadNewVersion": true}, nil)
//
//}
//
//func checkLatestVersion(ctx context.Context, urlPrefix string) (info *ReleaseInfo, err error) {
//	resp, err := gclient.New().Get(ctx, urlPrefix+ReleaseInfoJson)
//	if err != nil {
//		return nil, err
//	}
//
//	err = gconv.Scan(resp, &info)
//	if err != nil {
//		return nil, err
//	}
//
//	return info, err
//}
//
//func downloadNewVersion(ctx context.Context, urlPrefix string) error {
//
//	name := "RedisMan.exe.zip"
//
//	resp, err := gclient.New().Get(ctx, urlPrefix+name)
//	if err != nil {
//		return err
//	}
//
//	err = gcompress.UnZipContent(resp.ReadAll(), "./upgrade")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getCountry(ctx context.Context) (country string, err error) {
//	res, err := gclient.New().Get(ctx, "https://ipinfo.io/json")
//	if err != nil {
//		return "", err
//	}
//	type IpInfo struct {
//		Country string
//	}
//	var info IpInfo
//	err = gconv.Scan(res.ReadAllString(), &info)
//	return info.Country, err
//}
