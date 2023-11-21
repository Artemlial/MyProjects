package rrlb

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type LoadBalancer struct {
	Proxy    Proxy     `json:"proxy"`
	Backends []Backend `json:"backends"`
}

type Proxy struct {
	Port string `json:"port"`
}

type Backend struct {
	URL    string `json:"url"`
	IsDead atomic.Bool
}

func (b *Backend) SetDead(dead bool) {
	b.IsDead.Store(dead)
}

func (b *Backend) GetDead() bool {
	return b.IsDead.Load()
}

var lb LoadBalancer

func init() {
	f, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(f, &lb)
}

func Serve() {
	go healthCheck()

	s := http.Server{
		Addr:    lb.Proxy.Port,
		Handler: http.HandlerFunc(lbHandler),
	}

	log.Fatal(s.ListenAndServe())
}

func checkState(url *url.URL) bool {
	conn, err := net.DialTimeout("tcp", url.Host, time.Minute*1)
	if err != nil {
		log.Printf("dead backend on addr %s error: %s\n", url.Host, err.Error())
		return false
	}
	defer conn.Close()
	return true
}

func healthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			for _, back := range lb.Backends {
				pingURL, err := url.Parse(back.URL)
				if err != nil {
					log.Fatal(err)
				}
				lives := checkState(pingURL)
				back.SetDead(!lives)
				msg := "ok"
				if !lives {
					msg = "dead"
				}
				log.Printf("healthcheck addr: %s status: %s\n", pingURL.Host, msg)
			}
		}
	}
}

var mu sync.Mutex
var idx atomic.Uint32

func lbHandler(w http.ResponseWriter, r *http.Request) {
	mxlen := uint32(len(lb.Backends))

	mu.Lock()
	currBack := lb.Backends[idx.Load()%mxlen]
	for currBack.GetDead() {
		idx.Add(1)
		currBack = lb.Backends[idx.Load()%mxlen]
	}

	tgUrl, err := url.Parse(currBack.URL)
	if err != nil {
		log.Println(err.Error())
	}
	idx.Add(1)
	mu.Unlock()
	rp := httputil.NewSingleHostReverseProxy(tgUrl)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("%s is dead", tgUrl.Host)
		currBack.SetDead(true)
		lbHandler(w, r)
	}
	rp.ServeHTTP(w, r)
}
