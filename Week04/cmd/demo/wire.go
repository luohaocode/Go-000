//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/luohaocode/Go-000/Week04/ent"
	"github.com/luohaocode/Go-000/Week04/internal/biz"
	"github.com/luohaocode/Go-000/Week04/internal/data"
	"github.com/luohaocode/Go-000/Week04/internal/service"
)

func InitializeService(client *ent.Client) (*service.RecordUserService, error) {
	wire.Build(data.NewUserRepo, biz.NewRecordUserCase, service.NewRecordUserService)
	return &service.RecordUserService{}, nil
}