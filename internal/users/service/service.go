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
	Save(ctx context.Context, user model.User) (int, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Update(ctx context.Context, user model.User) error
}

type userServiceImpl struct {
	userRepository storage.UserStorage
}

func NewUserService(userRepo storage.UserStorage) UserService {
	return userServiceImpl{userRepository: userRepo}
}

func (u userServiceImpl) Update(ctx context.Context, user model.User) error {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "service_user_update")

	defer span.End()

	return u.userRepository.Update(childCtx, user)
}

func (u userServiceImpl) FindByID(ctx context.Context, id int) (model.User, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "service_user_find_by_id")

	defer span.End()

	return u.userRepository.FindByID(childCtx, id)
}

func (u userServiceImpl) Save(ctx context.Context, user model.User) (int, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "service_user_save")

	defer span.End()

	return u.userRepository.Save(childCtx, user)
}

func (u userServiceImpl) FindAll(ctx context.Context) ([]model.User, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "service_user_find_all")

	defer span.End()

	return u.userRepository.FindAll(childCtx)
}
