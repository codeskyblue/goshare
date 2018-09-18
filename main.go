package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/codeskyblue/comtool"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/levigross/grequests"
)

const PORT = 10351

var (
	version         = "0.1"
	ips             []*net.IP
	shares          map[string]string
	localServerAddr = fmt.Sprintf("http://127.0.0.1:%d", PORT)
)

func init() {
	shares = make(map[string]string)

	var err error
	ips, err = comtool.GetLocalIPs()
	if err != nil {
		log.Fatal(err)
	}
	if len(ips) <= 0 {
		log.Fatal("No network detected, can't share file to others.")
	}
}

func startServer() {
	m := mux.NewRouter()
	m.HandleFunc("/_version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(version))
	})

	m.HandleFunc("/_share", func(w http.ResponseWriter, r *http.Request) {
		path := r.FormValue("path")
		key := hashString(path)
		shares[key] = path
		w.Write([]byte(key))
	})

	m.HandleFunc("/_stop", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server quited"))
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()
	})
	m.HandleFunc("/{key}/{name}", func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		path := shares[key]
		http.ServeFile(w, r, path)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), handlers.LoggingHandler(os.Stdout, m))
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
	fServer     = kingpin.Flag("server", "run as server").Bool()
	fStopServer = kingpin.Flag("stop", "stop server").Bool()
)

func hashString(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil)[8:12])
}

func stopServer() {
	ro := &grequests.RequestOptions{}
	ro.RequestTimeout = 500 * time.Millisecond
	_, err := grequests.Get(localServerAddr+"/_stop", ro)
	if err == nil {
		log.Println("goshare stopped")
		return
	}
	log.Println("goshare already stopped")
	return
}

func wakeupServer() {
	if _, err := grequests.Get(localServerAddr+"/_version", nil); err == nil {
		return
	}
	log.Println("start goshare server")
	err := exec.Command(os.Args[0], "--server").Start()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *fServer {
		startServer()
		return
	}

	if *fStopServer {
		stopServer()
		return
	}

	wakeupServer() // if server not start, start it
	for _, name := range *sharePaths {
		if !comtool.Exists(name) {
			fmt.Printf("[%s] not exists\n", name)
			continue
		}
		fullpath, _ := filepath.Abs(name)
		basename := filepath.Base(name)
		ro := &grequests.RequestOptions{
			Params: map[string]string{"path": fullpath},
		}
		resp, err := grequests.Get(localServerAddr+"/_share", ro)
		if err != nil {
			log.Fatal(err)
		}
		key := resp.String()
		printLinks(key, basename)
	}
}
