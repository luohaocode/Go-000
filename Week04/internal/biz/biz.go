package biz

import (
	"context"
	"github.com/luohaocode/Go-000/Week04/ent"
)

type RecordRepo interface {
	SaveUser(context.Context, *ent.Client, *ent.User) (*ent.User, error)
}

func NewRecordUserCase(repo RecordRepo) *RecordUseCase {
	return &RecordUseCase{repo: repo}
}

type RecordUseCase struct {
	repo RecordRepo
}

func (uc *RecordUseCase) Record(ctx context.Context, client *ent.Client, u *ent.User) (*ent.User, error) {
	return uc.repo.SaveUser(ctx, client, u)
}
