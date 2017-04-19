package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/codeskyblue/comtool"
	"github.com/itang/gohttp"
)

var shareFolder string
var ips []*net.IP

const PORT = 10351

func init() {
	homedir, _ := comtool.HomeDir()
	shareFolder = filepath.Join(homedir, ".goshare")
	// println(shareFolder)
	if !comtool.Exists(shareFolder) {
		if err := os.Mkdir(shareFolder, 0755); err != nil {
			log.Fatal(err)
		}
	}

	var err error
	ips, err = comtool.GetLocalIPs()
	if err != nil {
		log.Fatal(err)
	}
	if len(ips) <= 0 {
		log.Fatal("No network detected, can't share file to others.")
	}
}

func startFileServer() {
	webroot := shareFolder
	server := &gohttp.FileServer{Port: PORT, Webroot: webroot}
	server.Start()
}

func printLinks(hash, name string) {
	for _, ip := range ips {
		if strings.Contains(ip.String(), ":") {
			continue // skip ipv6
		}
		url := fmt.Sprintf("http://%s:%d/%s/%s", ip, PORT, hash, name)
		fmt.Printf(comtool.Template(`{url}
wget {url} -O {name}
curl {url} -o {name}
`, map[string]string{
			"name": strconv.Quote(name),
			"url":  url,
		}))
		println("-------------------------------------")
	}
}

var (
	sharePaths  = kingpin.Arg("file", "shared file path").Strings()
	runAsServer = kingpin.Flag("server", "run as server").Bool()
)

func HashFile(fullpath string) string {
	h := md5.New()
	h.Write([]byte(fullpath))
	return fmt.Sprintf("%x", h.Sum(nil)[8:12])
}

func main() {
	kingpin.Parse()

	if *runAsServer {
		startFileServer()
		return
	}

	err := exec.Command(os.Args[0], "--server").Start()
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range *sharePaths {
		if comtool.Exists(name) {
			fullpath, _ := filepath.Abs(name)
			hash := HashFile(fullpath)

			basename := filepath.Base(name)
			sharepath := filepath.Join(shareFolder, hash, basename)
			if comtool.Exists(sharepath) {
				os.Remove(sharepath)
			}
			os.MkdirAll(filepath.Dir(sharepath), 0755)
			os.Symlink(fullpath, sharepath)
			printLinks(hash, basename)
		} else {
			fmt.Printf("[%s] not exists\n", name)
		}
	}
}
