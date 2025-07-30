package data

import (
	"context"
	"encoding/json"
	sqlc "music-streaming/sql/sqlc"
	"music-streaming/types"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgtype"
)

type SQLUserManagement interface {
	GetUser(ctx context.Context, username string) (*types.SubsonicUser, error)
	GetUsers(ctx context.Context) ([]*types.SubsonicUser, error)
	CreateUser(ctx context.Context, fields map[string]any) error
	UpdateUser(ctx context.Context, username string, fields map[string]string) error
	DeleteUser(ctx context.Context, username string) error
	ChangeUserPassword(ctx context.Context, username string, password string) error
}

func (d *DataLayerPg) GetUser(ctx context.Context, username string) (*types.SubsonicUser, error) {
	var cachedUser types.SubsonicUser
	userString, err := d.cache.Get(ctx, username).Result()
	if err != nil {
		conn, err := d.Pg_pool.Acquire(ctx)
		if err != nil {
			return nil, err
		}
		defer conn.Release()
		query := sqlc.New(conn)

		user, err := query.GetUserByUsername(ctx, pgtype.Text{String: username, Valid: true})
		if err != nil {
			return nil, &UserNotFoundError{username: username}
		}

		encodedUser, err := json.Marshal(types.MapSqlUserToSubsonicUser(&user, user.Password))
		if err != nil {
			return nil, err
		}
		if err = d.cache.Set(ctx, user.Username.String, encodedUser, time.Minute*10).Err(); err != nil {
			return nil, err
		}
		return types.MapSqlUserToSubsonicUser(&user, user.Password), nil
	}
	if err = json.Unmarshal([]byte(userString), &cachedUser); err != nil {
		return nil, err
	}
	return &cachedUser, nil
}

func (d *DataLayerPg) GetUsers(ctx context.Context) ([]*types.SubsonicUser, error) {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	sqlUsers, err := q.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*types.SubsonicUser, 0, len(sqlUsers))
	for _, user := range sqlUsers {
		temp := types.MapSqlUserToSubsonicUser(&user, "")
		temp.Password = ""
		users = append(users, temp)
	}
	return users, nil
}

func (d *DataLayerPg) CreateUser(ctx context.Context, fields map[string]any) error {
	insertSQL := goqu.Insert("users").Rows(fields)
	queryString, _, err := insertSQL.ToSQL()
	if err != nil {
		return err
	}

	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, queryString); err != nil {
		return err
	}
	return nil
}

func (d *DataLayerPg) UpdateUser(ctx context.Context, username string, fields map[string]string) error {
	sqlUpdate := goqu.Update("users").
		Set(fields).
		Where(goqu.Ex{"username": username})

	queryString, _, err := sqlUpdate.ToSQL()
	if err != nil {
		return err
	}

	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, queryString); err != nil {
		return err
	}
	return nil
}

func (d *DataLayerPg) DeleteUser(ctx context.Context, username string) error {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	if _, err = q.DeleteUser(ctx, pgtype.Text{String: username, Valid: true}); err != nil {
		return err
	}
	return nil
}

func (d *DataLayerPg) ChangeUserPassword(ctx context.Context, username string, password string) error {
	conn, err := d.Pg_pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	q := sqlc.New(conn)

	if _, err = q.ChangeUserPassword(ctx, sqlc.ChangeUserPasswordParams{Username: pgtype.Text{String: username, Valid: true}, Password: password}); err != nil {
		return err
	}
	return nil
}
