package storage

import (
	"context"
	"database/sql"
	"github.com/vielendanke/opentracing-example/internal/pkg/common"
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
		common.LogError(rowsErr)
		trace.AddSpanError(span, rowsErr)
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

	defer span.End()

	tx, txErr := u.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if txErr != nil {
		common.LogError(txErr)
		trace.AddSpanError(span, txErr)
		return 0, txErr
	}
	var id int
	row := tx.QueryRowContext(childCtx, "insert into users(username, email) values($1, $2) returning id",
		user.Username, user.Email)
	if scErr := row.Scan(&id); scErr != nil {
		common.LogError(scErr, tx.Rollback())
		trace.AddSpanError(span, txErr)
		return 0, scErr
	}
	common.LogError(tx.Commit())
	return id, nil
}

func (u sqlUserStorage) FindByID(ctx context.Context, id int) (model.User, error) {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "repository_user_save")

	defer span.End()

	var user model.User

	if scanErr := u.db.QueryRowContext(childCtx, "select u.id, u.username, u.email from users u where u.id=$1", id).
		Scan(&user.ID, &user.Username, &user.Email); scanErr != nil {
		common.LogError(scanErr)
		trace.AddSpanError(span, scanErr)
		return user, scanErr
	}
	return user, nil
}

func (u sqlUserStorage) Update(ctx context.Context, user model.User) error {
	childCtx, span := trace.NewSpanFromTracer(ctx, tracer, "repository_user_update")

	defer span.End()

	tx, txErr := u.db.BeginTx(childCtx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if txErr != nil {
		common.LogError(txErr)
		trace.AddSpanError(span, txErr)
		return txErr
	}
	if _, execErr := tx.Exec("update users set username=$1, email=$2 where id=$3", user.Username, user.Email, user.ID); execErr != nil {
		common.LogError(execErr)
		common.LogError(tx.Rollback())
		trace.AddSpanError(span, execErr)
		return execErr
	}
	return tx.Commit()
}
