//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ekbaya/asham/pkg/domain/services"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeServices(db *gorm.DB) (*services.ServiceContainer, error) {
	wire.Build(
		ServiceSet,
		services.NewServiceContainer,
	)
	return &services.ServiceContainer{}, nil
}
