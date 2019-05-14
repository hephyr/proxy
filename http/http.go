package http

import (
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"net/http"
)

type ProxyHandler struct {}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	if r.Method == http.MethodConnect {
		HTTPSProxy(w, r)
	} else {
		HTTPProxy(w, r)
	}
}

func HTTPProxy(w http.ResponseWriter, r *http.Request) {
	resp, err := request(r)
	if err != nil {
		log.Println(err)
		return
	}
	err = response(w, resp)
	if err != nil {
		log.Println(err)
	}
}

func request(r *http.Request) (*http.Response, error) {
	removeProxyHeaders(r)
	client := &http.Client{}
	return client.Do(r)
}

func response(w http.ResponseWriter, resp *http.Response) error {
	copyHeaders(w.Header(), resp.Header)
	_, err := io.Copy(w, resp.Body)
	return err
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

func HTTPSProxy(w http.ResponseWriter, r *http.Request) {
	client, err := HTTPSClient(w)
	if err != nil {
		return
	}

	server, err := HTTPSServer(r)
	if err != nil {
		return
	}

	handshake(server, client)

	go io.Copy(server, client)
	go io.Copy(client, server)
}

func HTTPSClient(w http.ResponseWriter) (net.Conn, error) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("HTTP Server does not support hijacking")
	}
	client, _, err := hij.Hijack()
	return client, err
}

func HTTPSServer(r *http.Request) (net.Conn, error) {
	host := r.URL.Host
	return net.Dial("tcp", host)
}

func handshake(server net.Conn, client net.Conn) {
	client.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))
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