package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	logger log.Logger
	Service
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &loggingMiddleware{logger, next}
	}
}

//func NewLoggingService(logger log.Logger, s Service) Service {
//	return &loggingMiddleware{logger, s}
//}

func (s *loggingMiddleware) Add(a, b int) (result int) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Add",
			"a", a,
			"b", b,
			"took", time.Since(begin),
			"result", result,
		)
	}(time.Now())
	return s.Service.Add(a, b)
}

func (s *loggingMiddleware) Subtract(a, b int) (result int) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Subtract",
			"a", a,
			"b", b,
			"took", time.Since(begin),
			"result", result,
		)
	}(time.Now())
	return s.Service.Subtract(a, b)
}

func (s *loggingMiddleware) Multiply(a, b int) (result int) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Multiply",
			"a", a,
			"b", b,
			"took", time.Since(begin),
			"result", result,
		)
	}(time.Now())
	return s.Service.Multiply(a, b)
}

func (s *loggingMiddleware) Divide(a, b int) (result int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Divide",
			"a", a,
			"b", b,
			"took", time.Since(begin),
			"result", result,
			"err", err,
		)
	}(time.Now())
	return s.Service.Divide(a, b)
}
