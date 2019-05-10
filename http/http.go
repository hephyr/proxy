package http

import (
	"io"
	"log"
	"net"
	"net/http"
)

type ProxyHandler struct {}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := request(r)
	if err != nil {
		log.Println(err)
		return
	}
	copyHeaders(w.Header(), resp.Header)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println(err)
	}
}

func request(r *http.Request) (*http.Response, error) {
	removeProxyHeaders(r)
	client := &http.Client{}
	return client.Do(r)
}

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func removeProxyHeaders(r *http.Request) {
	r.RequestURI = ""
}

func Listen(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	err = http.Serve(l, &ProxyHandler{})
	if err != nil {
		log.Fatal(err)
	}
}