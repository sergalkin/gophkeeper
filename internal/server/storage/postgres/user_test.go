package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/pkg/utils"
)

func TestNewPostgresUserStorage(t *testing.T) {
	tests := []struct {
		name string
		want *UserPostgresStorage
	}{
		{
			name: "New Postgres User Storage can be created",
			want: &UserPostgresStorage{conn: &pgx.Conn{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPostgresUserStorage(&pgx.Conn{}), "NewPostgresUserStorage()")
		})
	}
}

func TestUserPostgresStorage_Create(t *testing.T) {
	utils.RefreshTestDatabase()

	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name    string
		con     *pgx.Conn
		args    args
		wantErr bool
		do      func(user model.User, storage UserPostgresStorage)
	}{
		{
			name: "User can be created",
			con:  con,
			args: args{
				ctx:  ctx,
				user: model.User{Login: "test", Password: "test"},
			},
			wantErr: false,
			do: func(user model.User, storage UserPostgresStorage) {
			},
		},
		{
			name: "Conflict error will be thrown on attempt to create a user with same login",
			con:  con,
			args: args{
				ctx:  ctx,
				user: model.User{Login: "test", Password: "test"},
			},
			do: func(user model.User, storage UserPostgresStorage) {
				storage.Create(context.Background(), user)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserPostgresStorage{conn: tt.con}
			tt.do(tt.args.user, u)

			got, err := u.Create(tt.args.ctx, tt.args.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NotNil(t, got.ID)
		})
	}
}

func TestUserPostgresStorage_GetByLoginAndPassword(t *testing.T) {
	utils.RefreshTestDatabase()

	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name      string
		con       *pgx.Conn
		args      args
		want      model.User
		wantEmpty bool
		do        func(user model.User, storage UserPostgresStorage)
	}{
		{
			name: "User can be gotten",
			con:  con,
			args: args{
				ctx:  ctx,
				user: model.User{Login: "test", Password: "test"},
			},
			wantEmpty: false,
			do: func(user model.User, storage UserPostgresStorage) {
				storage.Create(context.Background(), user)
			},
		},
		{
			name: "User model will have nil ID on failed user retrieval",
			con:  con,
			args: args{
				ctx:  ctx,
				user: model.User{Login: "test", Password: "test"},
			},
			wantEmpty: true,
			do: func(user model.User, storage UserPostgresStorage) {
				q, _ := con.Query(ctx, "Delete from users where login = 'test'")
				q.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserPostgresStorage{conn: tt.con}
			tt.do(tt.args.user, u)

			got, err := u.GetByLoginAndPassword(tt.args.ctx, tt.args.user)
			if tt.wantEmpty == true {
				assert.Nil(t, got.ID)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got.ID)
		})
	}
}

func TestUserPostgresStorage_DeleteUser(t *testing.T) {
	utils.RefreshTestDatabase()

	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	tests := []struct {
		name    string
		con     *pgx.Conn
		want    model.User
		wantErr assert.ErrorAssertionFunc
		do      func(ctx context.Context, user model.User) model.User
	}{
		{
			name: "User can be deleted by provided user.Model",
			con:  con,
			want: model.User{},
			do: func(ctx context.Context, user model.User) model.User {
				con.QueryRow(
					ctx, "insert into users (login, password) values ('test', 'test') returning id",
				).Scan(&user.ID)

				return user
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := model.User{Login: "test", Password: "test"}
			u := UserPostgresStorage{conn: con}

			user = tt.do(ctx, user)

			got, err := u.DeleteUser(ctx, user)
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteUser(%v, %v)", ctx, user)) {
				return
			}

			assert.Equalf(t, tt.want, got, "DeleteUser(%v, %v)", ctx, user)
		})
	}
}
