package main

import (
	"context"
	"encoding/json"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHealthHandler(cs Service, logger kitlog.Logger, zipkinTracer *zipkin.Tracer) http.Handler {
	router := mux.NewRouter()

	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, kitzipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		zipkinServer,
	}

	router.Path("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		router.ServeHTTP(w, r)
	}))

	healthEndpoint := MakeHealthCheckEndpoint(cs)
	healthEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(healthEndpoint)
	router.Path("/health").Handler(kithttp.NewServer(
		healthEndpoint,
		decodeHealthCheckRequest,
		encodeArithmeticResponse,
		options...,
	))

	return router
}

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(cs Service, logger kitlog.Logger, zipkinTracer *zipkin.Tracer) http.Handler {
	r := mux.NewRouter()

	zipkinServer := kitzipkin.HTTPServerTrace(zipkinTracer, kitzipkin.Name("http-transport"))


	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		zipkinServer,
	}

	arithmeticEndpoint := MakeArithmeticEndpoint(cs)
	rateLimiter := rate.NewLimiter(rate.Every(time.Second+1), 100)
	arithmeticEndpoint = NewTokenBucketLimiterWithBuildIn(rateLimiter)(arithmeticEndpoint)
	arithmeticEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "calculate-endpoint")(arithmeticEndpoint)

	r.Path("/calculate/{type}/{a}/{b}").Handler(kithttp.NewServer(
		arithmeticEndpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...
	)).Methods("POST")

	r.Path("/health").Handler(kithttp.NewServer(
		MakeHealthCheckEndpoint(cs),
		decodeHealthCheckRequest,
		encodeArithmeticResponse,
		options...,
	))

	return r
}

// decodeArithmeticRequest decode request params to struct
func decodeArithmeticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	a, _ := strconv.Atoi(pa)
	b, _ := strconv.Atoi(pb)

	return ArithmeticRequest{
		RequestType: requestType,
		A:           a,
		B:           b,
	}, nil
}

// encodeArithmeticResponse encode response to return
func encodeArithmeticResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}
