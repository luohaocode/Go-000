package data

import (
	"context"
	"demo/internal/biz"
	"fmt"
	"github.com/luohaocode/Go-000/Week04/ent"
	"log"
)

var _ biz.RecordRepo = (*userRepo)(nil)

type userRepo struct{}

func NewUserRepo() biz.RecordRepo {
	return new(userRepo)
}

func (ur *userRepo) SaveUser(ctx context.Context, client *ent.Client, u *ent.User) (user *ent.User, err error) {
	user, err = client.User.
		Create().
		SetName(u.Name).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created: ", u)
	return user, err
}
