package version

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gbuild"
	"os"
)

func VersionAction() {
	if len(os.Args) == 2 && os.Args[1] == "-v" {
		fmt.Println(gbuild.Get("version"))
		fmt.Println(gbuild.Info())
		fmt.Println(gbuild.Data())
		panic(1)
		os.Exit(0)
	}
}
