package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Error          error
	handlers       []HandlerFunc
}

type HandlerFunc func(ctx *Request)

func NewRequest(w http.ResponseWriter, r *http.Request) *Request {
	return &Request{
		Request:        r,
		ResponseWriter: w,
	}
}

func (r *Request) Fail(status int) {
	message := http.StatusText(status)
	if r.Error != nil {
		message = r.Error.Error()
	} else {
		logrus.WithField("status", status).Warn("request failed without explicit error")
	}

	r.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.ResponseWriter.WriteHeader(status)
	if r.Request != nil && r.Request.Method == http.MethodHead {
		return
	}
	_, _ = r.ResponseWriter.Write([]byte(message))
}
