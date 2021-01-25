package middleware

import "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/endpoints"

type ServiceMiddleware func(endpoints.UserSerivcer) endpoints.UserSerivcer

func MakeServiceMiddleware(s endpoints.UserSerivcer) endpoints.UserSerivcer {
	mids := []ServiceMiddleware{
		makePrometheusMiddleware,
	}
	for _, m := range mids {
		s = m(s)
	}

	return s
}

