package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"

	userPb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/endpoints"
	"github.com/prometheus/client_golang/prometheus"
	kitPrometheus "github.com/go-kit/kit/metrics/prometheus"
)

var (
	fieldKeys = []string{"method", "error"}
)

func makePrometheusMiddleware(next endpoints.UserSerivcer) endpoints.UserSerivcer {
	requestCount := kitPrometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "usersvc_space",
		Subsystem: "usersvc",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "usersvc_space",
		Subsystem: "usersvc",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitPrometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "usersvc_space",
		Subsystem: "usersvc",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})
	return &prometheusMiddleware{requestCount, requestLatency, countResult, next}
}

type prometheusMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           endpoints.UserSerivcer
}

func (pm *prometheusMiddleware) timeDiff(method string, begin time.Time, err error) {
	lvs := []string{"method", method, "error", fmt.Sprint(err != nil)}
	pm.requestCount.With(lvs...).Add(1)
	pm.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
}

func (pm *prometheusMiddleware) GetUser(ctx context.Context, req string) (res *userPb.GetUserResponse, err error) {
	defer pm.timeDiff("GetUser", time.Now(), err)

	res, err = pm.next.GetUser(ctx, req)
	return
}

func (pm *prometheusMiddleware) Login(ctx context.Context, req endpoints.LoginRequest) (res *userPb.LoginResponse, err error) {
	defer pm.timeDiff("Login", time.Now(), err)

	res, err = pm.next.Login(ctx, req)
	return
}

func (pm *prometheusMiddleware) SendCode(ctx context.Context) (res *userPb.SendCodeResponse, err error) {
	defer pm.timeDiff("SendCode", time.Now(), err)

	res, err = pm.next.SendCode(ctx)
	return
}

func (pm *prometheusMiddleware) Register(ctx context.Context, req endpoints.RegisterRequest) (err error) {
	defer pm.timeDiff("Register", time.Now(), err)

	err = pm.next.Register(ctx, req)
	return
}

func (pm *prometheusMiddleware) UserList(ctx context.Context, req endpoints.UserListRequest) (res *userPb.
UserListResponse, err error) {
	defer pm.timeDiff("UserList", time.Now(), err)

	res, err = pm.next.UserList(ctx, req)
	return
}

func (pm *prometheusMiddleware) Logout(ctx context.Context, req endpoints.LogoutRequest) (err error) {
	defer pm.timeDiff("Logout", time.Now(), err)

	err = pm.next.Logout(ctx, req)
	return
}