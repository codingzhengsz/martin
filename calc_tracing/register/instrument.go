package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"golang.org/x/time/rate"
	"time"
)

var ErrLimitExceed = errors.New("Rate limit exceed! ")

// NewTokenBucketLimiterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

type instrumentService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentService) Add(a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Add"}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Add(a, b)
}

func (s *instrumentService) Multiply(a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Multiply"}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Multiply(a, b)
}

func (s *instrumentService) Subtract(a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Subtract"}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Subtract(a, b)
}

func (s *instrumentService) Divide(a, b int) (int, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Divide"}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Divide(a, b)
}

func (s *instrumentService) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		s.requestCount.With(lvs...).Add(1)
		s.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	result = s.Service.HealthCheck()
	return
}
