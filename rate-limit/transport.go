package main

import (
	"context"
	"encoding/json"
	"errors"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(cs Service, logger kitLog.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kitHttp.ServerOption{
		kitHttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kitHttp.ServerErrorEncoder(kitHttp.DefaultErrorEncoder),
	}

	arithmeticEndpoint := MakeArithmeticEndpoint(cs)
	rateLimiter := rate.NewLimiter(rate.Every(time.Second + 1), 3)
	arithmeticEndpoint = NewTokenBucketLimiterWithBuildIn(rateLimiter)(arithmeticEndpoint)

	r.Path("/calculate/{type}/{a}/{b}").Handler(kitHttp.NewServer(
		arithmeticEndpoint,
		decodeArithmeticRequest,
		encodeArithmeticResponse,
		options...
	)).Methods("POST")

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
