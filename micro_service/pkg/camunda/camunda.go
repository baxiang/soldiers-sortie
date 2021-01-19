package camunda

import (
	"context"
	"userService/pkg/camunda/pb"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
)

type Service struct {
	ProcessDefinition pb.ProcessDefinitionServer
	ProcessInstance   pb.ProcessInstanceServer
	Task              pb.TaskServer
	ExternalTask      pb.ExternalTaskServer
}

var defaultService *Service

func Load(client consul.Client, log log.Logger) {
	defaultService = new(Service)
	var (
		tags        []string
		passingOnly = true
	)
	instancer := consul.NewInstancer(client, log, "camundaService", tags, passingOnly)
	defaultService.ProcessDefinition = GetProcessDefinitionConsulEndpoints(instancer, log)
	defaultService.ProcessInstance = GetProcessInstanceConsulEndpoints(instancer, log)
	defaultService.Task = GetTaskConsulEndpoints(instancer, log)
	defaultService.ExternalTask = GetExternalTaskConsulEndpoints(instancer, log)
}

func Get() *Service {
	return defaultService
}

func encodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
