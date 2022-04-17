package storage

import (
	"context"
	"database/sql"
	"github.com/vielendanke/opentracing-example/internal/pkg/trace"
	"github.com/vielendanke/opentracing-example/internal/users/model"
	"log"
)

var tracer = trace.NewTracer("users_repository")

type CRUD interface {
	FindAll(ctx context.Context) ([]model.User, error)
	Save(ctx context.Context, user model.User) (int, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	Update(ctx context.Context, user model.User) error
}

type UserStorage interface {
	CRUD
}

type sqlUserStorage struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserStorage {
	return sqlUserStorage{db: db}
}

func (u sqlUserStorage) FindAll(ctx context.Context) ([]model.User, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "repository_user_find_all")

	defer span.End()

	rows, rowsErr := u.db.QueryContext(childCtx, "select u.id, u.username, u.email from users u")

	if rowsErr != nil {
		trace.AddSpanError(span, rowsErr)
		trace.FailSpan(span, "Database is not available")
		return nil, rowsErr
	}
	allUsers := make([]model.User, 0)

	for rows.Next() {
		var usr model.User
		if scErr := rows.Scan(&usr.ID, &usr.Username, &usr.Email); scErr != nil {
			log.Println(scErr)
			continue
		}
		allUsers = append(allUsers, usr)
	}
	return allUsers, nil
}

func (u sqlUserStorage) Save(ctx context.Context, user model.User) (int, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "repository_user_save")

	tx, txErr := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if txErr != nil {
		log.Println(txErr)
		trace.AddSpanError(span, txErr)
		trace.FailSpan(span, "Transaction failed to start")
		return 0, txErr
	}
	var id int
	row := tx.QueryRowContext(childCtx, "insert into users(username, email) values($1, $2) returning id",
		user.Username, user.Email)
	if scErr := row.Scan(&id); scErr != nil {
		log.Println(scErr)
		trace.AddSpanError(span, txErr)
		trace.FailSpan(span, "Failed to scan ID of user")
		return 0, scErr
	}
	return id, nil
}

func (u sqlUserStorage) FindByID(ctx context.Context, id int) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u sqlUserStorage) Update(ctx context.Context, user model.User) error {
	//TODO implement me
	panic("implement me")
}
