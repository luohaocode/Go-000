package main

import (
	"context"
	v1 "github.com/luohaocode/Go-000/Week04/api/demo/v1"
	grpcServer "github.com/luohaocode/Go-000/Week04/cmd/demo/server"
	"github.com/luohaocode/Go-000/Week04/config"
	"github.com/luohaocode/Go-000/Week04/internal/pkg"
	"log"
)

func main() {
	svr := grpcServer.NewServer(":9000")

	c := new(config.DbConfig)
	_ = config.ApplyYAML(config.GetYAML("config/demo.yaml"), c)
	dbClient, _ := config.New(c.Address, c.Name, c.Username, c.Password, config.ApplyOptions(c)...)
	conn, _ := config.Connect(dbClient)
	defer conn.Close()

	//rs := service.NewRecordUserService(conn, biz.NewRecordUserCase(data.NewUserRepo()))
	rs, _ := InitializeService(conn)

	v1.RegisterHelloServiceServer(svr.Server, rs)
	svr.Start(context.Background())

	app := pkg.New()
	if err := app.Run(); err != nil {
		log.Println("app failed:% v\n", err)
		return
	}
}
