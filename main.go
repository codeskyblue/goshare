// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codeskyblue/comtool"
	"github.com/itang/gohttp"
)

var shareFolder string
var ips []*net.IP

const PORT = 10351

func init() {
	homedir, _ := comtool.HomeDir()
	shareFolder = filepath.Join(homedir, ".goshare")
	println(shareFolder)
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

func printLinks(name string) {
	for _, ip := range ips {
		url := fmt.Sprintf("http://%s:%d/%s", ip, PORT, name)
		fmt.Printf(comtool.Template(`{url}
wget {url} -O "{name}"
curl -o "{name}" {url}
`, map[string]string{
			"name": name,
			"url":  url,
		}))
		println("-------------------------------------")
	}
}

func main() {
	fserver := flag.Bool("server", false, "start file server")
	flag.Usage = func() {
		fmt.Println(`Usage of goshare:
    goshare <file> : generate download links`)
	}

	flag.Parse()
	if *fserver {
		startFileServer()
	}
	err := exec.Command(os.Args[0], "-server").Start()
	fmt.Println(err)
	for _, name := range flag.Args() {
		if comtool.Exists(name) {
			fullpath, _ := filepath.Abs(name)
			basename := filepath.Base(name)
			sharepath := filepath.Join(shareFolder, basename)
			if comtool.Exists(sharepath) {
				os.Remove(sharepath)
			}
			os.Symlink(fullpath, sharepath)
			printLinks(basename)
		} else {
			fmt.Printf("[%s] not exists\n", name)
		}
	}
}
