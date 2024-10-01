package bootserver

import (
	"myproject/pkg/admin"
	"myproject/pkg/config"
	"myproject/pkg/user"
	"myproject/pkg/vendor"

	"github.com/labstack/echo/v4"
)

type ServerHttp struct {
	engine *echo.Echo
}

func NewServerHttp(userHandler user.Handler, vendorHandler vendor.Handler, adminHandler admin.Handler) *ServerHttp {
	engine := echo.New()

	// Mount user routes
	userHandler.MountRoutes(engine)

	// Mount vendor routes
	vendorHandler.MountRoutes(engine)
	adminHandler.MountRoutes(engine)
	//return &ServerHttp{Engine: engine}
	return &ServerHttp{engine}
}

func (s *ServerHttp) Start(conf config.Config) {
	s.engine.Start(conf.Host + ":" + conf.ServerPort)
}
