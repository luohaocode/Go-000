package main

import (
	"github.com/google/wire"
	v1 "github.com/luohaocode/Go-000/Week04/api/demo/v1"
	grpcServer "github.com/luohaocode/Go-000/Week04/cmd/demo/server"
	"github.com/luohaocode/Go-000/Week04/config"
	"github.com/luohaocode/Go-000/Week04/ent"
	"github.com/luohaocode/Go-000/Week04/internal/biz"
	"github.com/luohaocode/Go-000/Week04/internal/data"
	"github.com/luohaocode/Go-000/Week04/internal/pkg"
	"github.com/luohaocode/Go-000/Week04/internal/service"
	"log"
)

func main() {
	svr := grpcServer.NewServer("9000")

	c := new(config.DbConfig)
	_ = config.ApplyYAML(config.GetYAML("../../config/demo.yaml"), c)
	dbClient, _ := config.New(c.Address, c.Name, c.Username, c.Password, config.ApplyOptions(c)...)
	conn, _ := config.Connect(dbClient)
	defer conn.Close()

	//rs := service.NewHelloWorldService(conn, biz.NewRecordUserCase(data.NewUserRepo()))
	rs, _ := initializeHLService(conn)
	v1.RegisterHelloServiceServer(svr.Server, &rs)

	app := pkg.New()
	if err := app.Run(); err != nil {
		log.Println("app failed:% v\n", err)
		return
	}
}

func initializeHLService(client *ent.Client) (service.RecordUserService, error) {
	var serviceSet = wire.NewSet(data.NewUserRepo, biz.NewRecordUserCase, wire.Value(service.RecordUserService{Client: client}))
	wire.Build(serviceSet)
	return service.RecordUserService{}, nil
}
