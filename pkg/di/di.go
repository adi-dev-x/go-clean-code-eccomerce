package di

import (
	"myproject/pkg/admin"
	bootserver "myproject/pkg/boot"
	"myproject/pkg/vendor"

	services "myproject/pkg/client"
	"myproject/pkg/config"
	db "myproject/pkg/database"
	"myproject/pkg/user"
)

func InitializeEvent(conf config.Config) (*bootserver.ServerHttp, error) {

	sqlDB, err := db.ConnectPGDB(conf)
	if err != nil {
		return nil, err // Return early if there's an error connecting to the database
	}

	// Create a new repository using the *sql.DB instance
	userRepository := user.NewRepository(sqlDB)
	userService := user.NewService(userRepository)
	myService := services.MyService{Config: conf}
	// admjwt := middleware.Adminjwt{Config: conf}
	admjwt := user.Adminjwt{Config: conf}
	userHandler := user.NewHandler(userService, myService, admjwt)

	vendorRepository := vendor.NewRepository(sqlDB)
	vendorService := vendor.NewService(vendorRepository)
	vendorjwt := vendor.Vendorjwt{Config: conf} // Corrected import path
	vendorHandler := vendor.NewHandler(vendorService, myService, vendorjwt)

	adminRepository := admin.NewRepository(sqlDB)
	adminService := admin.NewService(adminRepository)
	// adminjwt := admin.adminjwt{Config: conf} // Corrected import path
	adminjwt := admin.Adminjwt{Config: conf}
	adminHandler := admin.NewHandler(adminService, myService, adminjwt)

	serverHttp := bootserver.NewServerHttp(*userHandler, *vendorHandler, *adminHandler)

	return serverHttp, nil
}
