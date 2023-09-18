package model

type Response struct {
	ID      int64
	Code    int
	Message string
	Body    string
	Headers map[string][]string
}
