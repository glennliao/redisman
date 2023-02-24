package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var configTemp = `
gfcli:
  build:
    name:     "redisman"
    arch:     "arm64,amd64"
    system:   "windows,linux,darwin"
    output:   "./temp"
    extra:    "-trimpath"
    varMap:
      version: {{version}}
`

var version = ""

var buildDir = "./build"

func main() {

	if len(os.Args) > 1 {
		if os.Args[1] == "zip" {
			gcompress.ZipPath(os.Args[2], filepath.Join(os.Args[3], fmt.Sprintf("redisman_desktop_%s_%s.zip", runtime.GOOS, runtime.GOARCH)))
			return
		}
	}

	buildDir, _ = filepath.Abs(buildDir)

	gfile.Mkdir(buildDir)

	out, err := gproc.ShellExec(nil, "git describe --abbrev=0 --tags")
	if err != nil {
		panic(err)
	}
	version = strings.TrimSpace(out)
	log.Println("build redisman server :" + version)
	server()
	log.Println("build redisman desktop :" + version)
	//desktop()

	gfile.PutContents(filepath.Join(buildDir, "release_info.json"), gjson.New(g.Map{
		"Version": version,
	}).MustToJsonString())

	repository := os.Getenv("GITHUB_REPOSITORY")

	url := "https://api.github.com/repos/" + repository + "/commits/" + version
	g.Log().Info(nil, url)

	res, err := g.Client().Get(nil, url)
	handleError(err)

	msg := gjson.New(res.ReadAllString()).Get("commit.message").String()
	if msg == "" {
		msg = "can't get message"
	}

	gfile.PutContents(filepath.Join("", "./temp/changelog.md"), msg)
}

func server() {
	os.Chdir("./cmd/server")
	configTemp = strings.ReplaceAll(configTemp, "{{version}}", version)
	gfile.PutContents("./config.yaml", configTemp)
	out, err := gproc.ShellExec(nil, "gf build")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	list, err := gfile.Glob("./temp/*/*")
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		newPath := filepath.Join(buildDir, "redisman_server_"+filepath.Base(filepath.Dir(item)))
		if strings.Contains(item, "linux") {
			newPath += ".tar.gz"
			TarGz(filepath.Join("", item), newPath)
		} else {
			newPath += ".zip"
			gcompress.ZipPath(filepath.Join("", item), newPath)
		}

	}

	os.Chdir("../../")
}

func desktop() {
	os.Chdir("./cmd/desktop")
	out, err := gproc.ShellExec(nil, "wails build -trimpath ")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	list, err := gfile.Glob("./build/bin/*")
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		gcompress.ZipPath(item, filepath.Join(buildDir, fmt.Sprintf("redisman_desktop_%s_%s.zip", runtime.GOOS, runtime.GOARCH)))

	}

	os.Chdir("../../")
}
func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func TarGz(srcDirPath string, destFilePath string) {
	fw, err := os.Create(destFilePath)
	handleError(err)
	defer fw.Close()

	// Gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// Tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Check if it's a file or a directory
	f, err := os.Open(srcDirPath)
	handleError(err)
	fi, err := f.Stat()
	handleError(err)
	if fi.IsDir() {
		// handle source directory
		fmt.Println("Cerating tar.gz from directory...")
		tarGzDir(srcDirPath, srcDirPath, tw)
	} else {
		// handle file directly
		fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
		tarGzFile(srcDirPath, fi.Name(), tw, fi)
	}
	fmt.Println("Well done!")
}

func tarGzDir(srcDirPath string, recPath string, tw *tar.Writer) {
	// Open source diretory
	dir, err := os.Open(srcDirPath)
	handleError(err)
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	handleError(err)
	for _, fi := range fis {
		// Append path
		curPath := srcDirPath + "/" + fi.Name()
		// Check it is directory or file
		if fi.IsDir() {
			// Directory
			// (Directory won't add unitl all subfiles are added)
			fmt.Printf("Adding path...%s\\n", curPath)
			tarGzDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			fmt.Printf("Adding file...%s\\n", curPath)
		}

		tarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
}

// Deal with files
func tarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) {
	if fi.IsDir() {
		// Create tar header
		hdr := new(tar.Header)
		// if last character of header name is '/' it also can be directory
		// but if you don't set Typeflag, error will occur when you untargz
		hdr.Name = recPath + "/"
		hdr.Typeflag = tar.TypeDir
		hdr.Size = 0
		//hdr.Mode = 0755 | c_ISDIR
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err := tw.WriteHeader(hdr)
		handleError(err)
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		handleError(err)
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = filepath.Base(recPath)
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		handleError(err)

		// Write file data
		_, err = io.Copy(tw, fr)
		handleError(err)
	}
}
