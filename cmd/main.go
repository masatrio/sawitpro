package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/handler"
	"github.com/sawitpro/UserService/helper/hasher/bcrypt"
	"github.com/sawitpro/UserService/repository/postgres"
	"github.com/sawitpro/UserService/service/service"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)

	mw, err := handler.InitMiddleware()
	if err != nil {
		panic(fmt.Sprintf("error creating middleware, err = %s", err.Error()))
	}

	e.Use(mw...)

	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() generated.ServerInterface {
	dbDsn, ok := os.LookupEnv("DATABASE_URL")

	if !ok {
		panic("DATABASE_URL env not set")
	}

	TZ, ok := os.LookupEnv("TZ")
	if !ok {
		TZ = "Asia/Bangkok"
	}
	os.Setenv("TZ", TZ)
	time.LoadLocation(TZ)

	repository := postgres.NewClient(postgres.ClientOptions{
		DSN: dbDsn,
	})

	svc := service.NewService(service.ServiceOpts{
		Repository: repository,
		Hasher:     bcrypt.NewHasher(),
	})

	opts := handler.ServerOpts{
		Service: svc,
	}
	return handler.NewServer(opts)
}
