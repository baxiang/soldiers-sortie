package gateway

import "time"

var (
	rpcRetryTimes = 1
	rpcTimeOut    = 5 * time.Second
)

const (
	scanBreaker        = "scanBreaker"
	staticBreaker      = "staticBreaker"
	institutionBreaker = "institutionBreaker"
	userBreaker        = "userBreaker"
	merchantBreaker    = "merchantBreaker"
	termBreaker        = "termBreaker"
	workflowBreaker    = "workflowBeaker"
)

func SetRetryTime(times int) {
	rpcRetryTimes = times
}

func SetTimeOut(timeOut time.Duration) {
	rpcTimeOut = timeOut
}
