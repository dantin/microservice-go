package service

import "net/http"

// Route represents a route endpoint.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a collection of `Route`.
type Routes []Route

var routes = Routes{
	Route{"index", "GET", "/", index},
	Route{"getAccount", "GET", "/accounts/{accountId}", getAccount},
}
