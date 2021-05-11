package internalhttp

import (
	"net/http"
)

type Router interface {
	Routes() []Route
}

type Route struct {
	Name   string
	Method string
	Path   string
	Func   http.HandlerFunc
}
