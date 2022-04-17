package service

import (
	"context"
	"github.com/vielendanke/opentracing-example/internal/pkg/trace"
	"github.com/vielendanke/opentracing-example/internal/users/model"
	"github.com/vielendanke/opentracing-example/internal/users/storage"
)

var tracer = trace.NewTracer("user_service")

type UserService interface {
	FindAll(ctx context.Context) ([]model.User, error)
}

type userServiceImpl struct {
	userRepository storage.UserStorage
}

func NewUserService(userRepo storage.UserStorage) UserService {
	return userServiceImpl{userRepository: userRepo}
}

func (u userServiceImpl) FindAll(ctx context.Context) ([]model.User, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "service_user_find_all")

	defer span.End()

	return u.userRepository.FindAll(childCtx)
}
