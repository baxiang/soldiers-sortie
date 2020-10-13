// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sdgmf/go-project-sample/internal/app/ratings"
	"github.com/sdgmf/go-project-sample/internal/app/ratings/controllers"
	"github.com/sdgmf/go-project-sample/internal/app/ratings/grpcservers"
	"github.com/sdgmf/go-project-sample/internal/app/ratings/services"
	"github.com/sdgmf/go-project-sample/internal/app/ratings/repositories"
	"github.com/sdgmf/go-project-sample/internal/pkg/app"
	"github.com/sdgmf/go-project-sample/internal/pkg/config"
	"github.com/sdgmf/go-project-sample/internal/pkg/consul"
	"github.com/sdgmf/go-project-sample/internal/pkg/jaeger"

	"github.com/sdgmf/go-project-sample/internal/pkg/database"
	"github.com/sdgmf/go-project-sample/internal/pkg/log"
	"github.com/sdgmf/go-project-sample/internal/pkg/transports/grpc"
	"github.com/sdgmf/go-project-sample/internal/pkg/transports/http"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	database.ProviderSet,
	services.ProviderSet,
	consul.ProviderSet,
	jaeger.ProviderSet,
	http.ProviderSet,
	grpc.ProviderSet,
	ratings.ProviderSet,
	repositories.ProviderSet,
	controllers.ProviderSet,
	grpcservers.ProviderSet,
)


func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
