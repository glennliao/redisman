package action

//
//func init() {
//	ws.RegAction("relaunch", relaunch)
//}
//
//const UpgradeTempWaitDelFile = "_tmp_wait_del"
//
//func relaunch(ctx context.Context, req *ws.Req, reply func(ctx context.Context, ret any, err error)) {
//	err := gfile.Rename(os.Args[0], UpgradeTempWaitDelFile)
//	if err != nil {
//		reply(ctx, nil, err)
//		return
//	}
//
//	err = gfile.Remove(os.Args[0])
//	if err != nil {
//		reply(ctx, nil, err)
//		return
//	}
//
//	cmd := exec.Command(os.Args[0])
//	cmd.Run()
//}
