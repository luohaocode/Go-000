package service

import (
	"context"
	"fmt"
	v1 "github.com/luohaocode/Go-000/Week04/api/demo/v1"
	"github.com/luohaocode/Go-000/Week04/ent"
	"github.com/luohaocode/Go-000/Week04/internal/biz"
)

type RecordUserService struct {
	v1.UnimplementedHelloServiceServer
	ruc    *biz.RecordUseCase
	Client *ent.Client
}

func NewHelloWorldService(client *ent.Client, ruc *biz.RecordUseCase) *RecordUserService {
	return &RecordUserService{
		ruc:    ruc,
		Client: client,
	}
}

func (s *RecordUserService) Record(req *v1.HelloRequest) (*v1.HelloResponse, error) {
	u := new(ent.User)
	u.Name = req.Name
	_, err := s.ruc.Record(context.Background(), s.Client, u)
	if err != nil {
		return &v1.HelloResponse{
			Message: fmt.Sprintf("Record user %s error", u.Name),
		}, err
	}

	return &v1.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
