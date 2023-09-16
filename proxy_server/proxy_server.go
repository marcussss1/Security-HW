package proxy_server

import (
	"crypto/tls"
	"errors"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"net/url"
)

type proxy struct{}

func StartServer() {
	e := echo.New()
	proxyHandler := &proxy{}
	e.Any("/*", proxyHandler.Handle)
	e.Logger.Fatal(e.Start(":8081"))
}

func StartServerTLS() {
	e := echo.New()
	proxyHandler := &proxy{}
	e.Any("/*", proxyHandler.Handle)
	
	server := &http.Server{
		Addr:         ":8081",
		Handler:      e,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	err := http2.ConfigureServer(server, nil)
	if err != nil {
		e.Logger.Fatal(err)
	}

	err = server.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func (p *proxy) Handle(ctx echo.Context) error {
	if ctx.Request().Method == http.MethodConnect {
		return p.HttpsHandle(ctx)
	}

	return p.HttpHandle(ctx)
}

func (p *proxy) HttpHandle(ctx echo.Context) error {
	requestUrlString := getUrlFromContext(ctx)
	requestUrl, err := url.Parse(requestUrlString)
	if err != nil {
		return err
	}

	ctx.Request().URL = requestUrl
	ctx.Request().Header.Del("Proxy-Connection")

	resp, err := http.DefaultTransport.RoundTrip(ctx.Request())
	if err != nil {
		return err
	}
	copyHeaders(ctx.Response().Header(), resp.Header)

	return ctx.NoContent(http.StatusOK)
}

func (p *proxy) HttpsHandle(ctx echo.Context) error {
	dest, err := net.Dial("tcp", ctx.Request().Host)
	if err != nil {
		return err
	}

	hijacker, ok := ctx.Response().Writer.(http.Hijacker)
	if !ok {
		return errors.New("hjacker")
	}

	client, _, err := hijacker.Hijack()
	if err != nil {
		return err
	}

	go transfer(dest, client)
	go transfer(client, dest)

	return nil
}

func getUrlFromContext(ctx echo.Context) string {
	protocol := ctx.Request().URL.Scheme
	host := ctx.Request().Host
	path := ctx.Request().URL.Path
	return protocol + "://" + host + path
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}