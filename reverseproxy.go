package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

const HOSTPORT = 80

type Host struct {
	HostAddress   string
	TargetAddress string
	Name          string
}

func (h *Host) getTargetURL() *url.URL {
	url, err := url.Parse("http://localhost:" + h.TargetAddress)
	check(err)
	return url
}

func (h *Host) getHostURL() *url.URL {
	url, err := url.Parse(h.HostAddress)
	check(err)
	return url
}

type Config struct {
	Hosts []Host
}

type HandlerFunc func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func addReversePortProxy(hostAddress *url.URL, targetAddress *url.URL, httpHandler HandlerFunc) {
	proxy := httputil.NewSingleHostReverseProxy(targetAddress)
	http.HandleFunc(hostAddress.String(), httpHandler(proxy))
	log.Printf("Redirecting %s to %s on port %s\n", hostAddress.String(), targetAddress.String(), strconv.Itoa(HOSTPORT))
}

func getConfig(configFile string) Config {
	f, err := os.Open(configFile)
	check(err)
	cfg := Config{}
	json.NewDecoder(bufio.NewReader(f)).Decode(&cfg)
	return cfg
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		p.ServeHTTP(w, r)
	}
}

func main() {
	cfg := getConfig("proxyconfig.json")
	for _, host := range cfg.Hosts {
		addReversePortProxy(host.getHostURL(), host.getTargetURL(), handler)
	}
	err := http.ListenAndServe(":"+strconv.Itoa(HOSTPORT), nil)
	check(err)
}
