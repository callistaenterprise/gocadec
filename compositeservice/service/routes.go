package service

import "net/http"

/**
 * Derived from http://thenewstack.io/make-a-restful-json-api-go/
 */
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"GetAccount",
		"GET",
		"/account/{accountId}",
		GetAccount,
	},
	Route{
		"Health",
		"GET",
		"/health",
		HealtCheck,
	},
}
