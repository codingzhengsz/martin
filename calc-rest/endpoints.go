package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"strings"
)

var (
	//ErrInvalidRequestType = fmt.Sprint()
	ErrInvalidRequestType = errors.New("RequestType has only four type: Add,Subtract,Multiply,Divide")
)

// ArithmeticRequest define request struct
type ArithmeticRequest struct {
	RequestType string `json:"request_type"`
	A           int    `json:"a"`
	B           int    `json:"b"`
}

// ArithmeticResponse define response struct
type ArithmeticResponse struct {
	Result  int    `json:"result"`
	Message string `json:"message"`
}

// CalculateEndpoint define endpoint
type ArithmeticEndpoint endpoint.Endpoint

// MakeArithmeticEndpoint make endpoint
func MakeArithmeticEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ArithmeticRequest)

		var (
			res, a, b int
			calError  error
		)

		a = req.A
		b = req.B

		if strings.EqualFold(req.RequestType, "Add") {
			res = svc.Add(a, b)
		} else if strings.EqualFold(req.RequestType, "Subtract") {
			res = svc.Subtract(a, b)
		} else if strings.EqualFold(req.RequestType, "Multiply") {
			res = svc.Multiply(a, b)
		} else if strings.EqualFold(req.RequestType, "Divide") {
			if res, calError = svc.Divide(a, b); calError != nil {
				return ArithmeticResponse{Result: res, Message: calError.Error()}, nil
			}
		} else {
			return nil, ErrInvalidRequestType
		}
		return ArithmeticResponse{Result: res, Message: ""}, nil
	}
}
