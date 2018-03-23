package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/fatih/color"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

type myTransport struct {
}

type Montioringpath struct {
	path        string
	count       int64
	duration    int64
	averageTime float64
}

var globalMap = make(map[string]Montioringpath)

func main() {
	colorize(color.FgGreen, "=> Reading Config")

	const (
		defaultPort        = "9090"
		defaultPortUsage   = "default server port"
		defaultTarget      = "http://127.0.0.1:8080"
		defaultTargetUsage = "default target port"
	)

	port := flag.String("port", defaultPort, defaultPortUsage)
	targetURL := flag.String("targetURL", defaultTarget, defaultTargetUsage)

	flag.Parse()

	colorize(color.FgCyan, "Server: ", *port)
	colorize(color.FgCyan, "Redirecting to: ", *targetURL)

	proxy := NewProxy(*targetURL)
	//http.HandleFunc("/proxyServer", ProxyServer)

	http.HandleFunc("/", proxy.handle)
	//log.Fatal(http.ListenAndServe(":"+*port, nil))
	log.Fatal(http.ListenAndServeTLS(":"+*port, "server.crt", "server.key", nil))
}

func NewProxy(target string) *Prox {
	url, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	return &Prox{
		target: url,
		proxy:  httputil.NewSingleHostReverseProxy(url),
	}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-SpaceHub", "SpaceHub")
	p.proxy.Transport = &myTransport{}
	p.proxy.ServeHTTP(w, r)
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()

	if request.Body != nil {
		buf, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		requestBodyInternal := ioutil.NopCloser(bytes.NewBuffer(buf))
		requestBody := ioutil.NopCloser(bytes.NewBuffer(buf))

		colorize(color.FgGreen, "=> Request Body: ", requestBodyInternal)
		request.Body = requestBody
	}

	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		colorize(color.FgRed, "=> Error response: ", err)
		return nil, err
	}

	elapsed := time.Since(start)
	key := request.Method + "-" + request.URL.Path

	val, ok := globalMap[key]
	if ok != true {
		var m Montioringpath
		m.path = request.URL.Path
		m.count = 0
		m.duration = 0
		m.averageTime = 0
		val = m
	}

	val.count++
	val.duration += elapsed.Nanoseconds()
	val.averageTime = float64(val.duration) / float64(val.count)
	globalMap[key] = val

	b, err := json.MarshalIndent(globalMap, "", "  ")
	if err != nil {
		colorize(color.FgRed, "=> Error: ", err)
	}
	colorize(color.FgCyan, b)

	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		colorize(color.FgRed, "=> Error in dump response: ", err)
		return nil, err
	}

	colorize(color.FgGreen, "=> Response Body: ", string(body))
	colorize(color.FgGreen, "=> Response Time: ", elapsed.Nanoseconds())

	return resp, nil
}
