package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		IndexHandler,
	},
	Route{
		"GetSamples",
		"GET",
		"/samples",
		GetSamplesHandler,
	},
	Route{
		"GetOneSample",
		"GET",
		"/samples/{sampleId}",
		GetOneSampleHandler,
	},
	Route{
		"PostSample",
		"POST",
		"/samples",
		PostSampleHandler,
	},
}
