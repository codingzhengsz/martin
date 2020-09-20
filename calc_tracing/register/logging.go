package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Add(a, b int) (result int) {
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

func (s *loggingService) Subtract(a, b int) (result int) {
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

func (s *loggingService) Multiply(a, b int) (result int) {
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

func (s *loggingService) Divide(a, b int) (result int, err error) {
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

func (s *loggingService) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		s.logger.Log(
			"function", "HealthCheck",
			"result", result,
			"took", time.Since(begin))
	}(time.Now())
	result = s.Service.HealthCheck()
	return
}
