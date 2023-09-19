package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"security/model"
	"security/parser"
	"security/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) Handler {
	return Handler{store: store}
}

func (h *Handler) GetRequests(ctx echo.Context) error {
	requests := h.store.GetRequests()

	return ctx.JSON(http.StatusOK, requests)
}

func (h *Handler) GetRequestByID(ctx echo.Context) error {
	request := h.store.GetRequestByID(ctx.Param("id"))

	return ctx.JSON(http.StatusOK, request)
}

func (h *Handler) RepeatRequestByID(ctx echo.Context) error {
	request := h.store.GetRequestByID(ctx.Param("id"))
	repeatRequest := parser.ParseRepeatRequest(request)

	ctx.Request().URL = repeatRequest.URL
	ctx.Request().Header.Del("Proxy-Connection")

	resp, err := http.DefaultTransport.RoundTrip(ctx.Request())
	if err != nil {
		panic(err)
	}
	copyHeaders(ctx.Response().Header(), resp.Header)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(responseBody)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// ScanRequestByID SQL INJECTION 2 ВАРИАНТ
func (h *Handler) ScanRequestByID(ctx echo.Context) error {
	request := h.store.GetRequestByID(ctx.Param("id"))
	repeatRequest := parser.ParseRepeatRequest(request)

	ctx.Request().URL = repeatRequest.URL
	ctx.Request().Header.Del("Proxy-Connection")

	resp, err := http.DefaultTransport.RoundTrip(ctx.Request())
	if err != nil {
		panic(err)
	}
	copyHeaders(ctx.Response().Header(), resp.Header)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, err = ctx.Response().Write(responseBody)
	if err != nil {
		panic(err)
	}

	if vulnerable := sendGetParamsVulnerableRequest(`'`, request, resp) ||
		sendGetParamsVulnerableRequest(`"`, request, resp); !vulnerable {
		fmt.Println("GET PARAMS УЯЗВИМЫ")
	}
	if vulnerable := sendPostParamsVulnerableRequest(`'`, request, resp) ||
		sendPostParamsVulnerableRequest(`"`, request, resp); !vulnerable {
		fmt.Println("POST PARAMS УЯЗВИМЫ")
	}
	if vulnerable := sendCookieParamsVulnerableRequest(`'`, request, resp) ||
		sendCookieParamsVulnerableRequest(`"`, request, resp); !vulnerable {
		fmt.Println("COOKIE PARAMS УЯЗВИМЫ")
	}
	if vulnerable := sendHTTPParamsVulnerableRequest(`'`, request, resp) ||
		sendHTTPParamsVulnerableRequest(`"`, request, resp); !vulnerable {
		fmt.Println("HTTP PARAMS УЯЗВИМЫ")
	}

	return nil
}

func sendGetParamsVulnerableRequest(char string, request model.Request, firstResponse *http.Response) bool {
	if len(request.GetParams) == 0 {
		return true
	}

	for k := range request.GetParams {
		for kk := range request.GetParams[k] {
			request.GetParams[k][kk] += char
		}
	}

	return sendVulnerableRequest(request, firstResponse)
}

func sendPostParamsVulnerableRequest(char string, request model.Request, firstResponse *http.Response) bool {
	if len(request.PostParams) == 0 {
		return true
	}

	for k := range request.PostParams {
		request.PostParams[k] += char
	}

	return sendVulnerableRequest(request, firstResponse)
}

func sendCookieParamsVulnerableRequest(char string, request model.Request, firstResponse *http.Response) bool {
	if len(request.Cookies) == 0 {
		return true
	}

	for k := range request.Cookies {
		request.Cookies[k] += char
	}

	return sendVulnerableRequest(request, firstResponse)
}

func sendHTTPParamsVulnerableRequest(char string, request model.Request, firstResponse *http.Response) bool {
	if len(request.Headers) == 0 {
		return true
	}

	for k := range request.Headers {
		for kk := range request.Headers[k] {
			request.Headers[k][kk] += char
		}
	}

	return sendVulnerableRequest(request, firstResponse)
}

func sendVulnerableRequest(request model.Request, firstResponse *http.Response) bool {
	parsedRequest := parser.ParseRepeatRequest(request)

	resp, err := http.DefaultTransport.RoundTrip(parsedRequest)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != firstResponse.StatusCode || resp.ContentLength != firstResponse.ContentLength {
		return true
	}

	return false
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
