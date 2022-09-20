package postgres

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/pkg/utils"
)

func TestNewSecretPostgresStorage(t *testing.T) {
	tests := []struct {
		name string
		want *SecretPostgresStorage
	}{
		{
			name: "New Postgres Secret Storage could be created",
			want: &SecretPostgresStorage{&pgx.Conn{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSecretPostgresStorage(&pgx.Conn{})

			assert.Equalf(t, tt.want, got, "NewPostgresSecretStorage() = %v, want %v", got, tt.want)
		})
	}
}

func TestSecretPostgresStorage_CreateSecret(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	uid := uuid.New()

	type args struct {
		ctx    context.Context
		secret model.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    model.Secret
		wantErr assert.ErrorAssertionFunc
		do      func()
	}{
		{
			name: "Secret can be created",
			args: args{
				ctx: ctx,
				secret: model.Secret{
					UserID:    uid,
					TypeID:    1,
					Title:     "Test",
					Content:   nil,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: nil,
				},
			},
			want: model.Secret{
				ID:        1,
				UserID:    uid,
				TypeID:    1,
				Title:     "Test",
				Content:   nil,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: nil,
			},
			wantErr: assert.NoError,
			do: func() {
				utils.RefreshTestDatabase()

				row, _ := con.Query(
					ctx, "insert into users (id, login, password) values ($1,$2,$3)", uid, "test", "test",
				)
				defer row.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()

			s := &SecretPostgresStorage{
				conn: con,
			}

			got, err := s.CreateSecret(tt.args.ctx, tt.args.secret)
			tt.wantErr(t, err, fmt.Sprintf("CreateSecret(%v, %v)", tt.args.ctx, tt.args.secret))

			assert.Equalf(t, tt.want, got, "CreateSecret(%v, %v)", tt.args.ctx, tt.args.secret)
		})
	}
}

func TestSecretPostgresStorage_GetSecret(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	uid := uuid.New()

	type args struct {
		ctx    context.Context
		secret model.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    model.Secret
		wantErr assert.ErrorAssertionFunc
		do      func()
	}{
		{
			name: "Secret can be gotten from database",
			args: args{
				ctx:    ctx,
				secret: model.Secret{ID: 1, UserID: uid},
			},
			wantErr: assert.NoError,
			do: func() {
				utils.RefreshTestDatabase()

				row, _ := con.Query(
					ctx, "insert into users (id, login, password) values ($1,$2,$3)", uid, "test", "test",
				)
				row.Close()

				r2, _ := con.Query(
					ctx, "insert into secrets (user_id, title, content, type_id) values ($1,$2,$3,$4)",
					uid, "test", hex.EncodeToString([]byte{10, 20}), 1,
				)
				r2.Close()
			},
		},
		{
			name: "Error will be returned when getting non existent secret",
			args: args{
				ctx:    ctx,
				secret: model.Secret{ID: 1, UserID: uid},
			},
			wantErr: assert.Error,
			do: func() {
				utils.RefreshTestDatabase()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()

			s := &SecretPostgresStorage{conn: con}

			got, err := s.GetSecret(tt.args.ctx, tt.args.secret)

			if !tt.wantErr(t, err, fmt.Sprintf("GetSecret(%v, %v)", tt.args.ctx, tt.args.secret)) {
				assert.NotNil(t, got.Content)
			}
		})
	}
}

func TestSecretPostgresStorage_DeleteSecret(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	uid := uuid.New()

	type args struct {
		ctx    context.Context
		secret model.Secret
	}
	tests := []struct {
		name    string
		args    args
		want    model.Secret
		wantErr assert.ErrorAssertionFunc
		do      func()
	}{
		{
			name: "Secret can be deleted",
			args: args{
				ctx:    ctx,
				secret: model.Secret{ID: 1, UserID: uid},
			},
			want:    model.Secret{},
			wantErr: assert.NoError,
			do: func() {
				utils.RefreshTestDatabase()

				row, _ := con.Query(
					ctx, "insert into users (id, login, password) values ($1,$2,$3)", uid, "test", "test",
				)
				row.Close()

				r2, _ := con.Query(
					ctx, "insert into secrets (user_id, title, content, type_id) values ($1,$2,$3,$4)",
					uid, "test", hex.EncodeToString([]byte{10, 20}), 1,
				)
				r2.Close()
			},
		},
		{
			name: "Error will be throws on deleting non existing  secret",
			args: args{
				ctx:    ctx,
				secret: model.Secret{ID: 1, UserID: uid},
			},
			want:    model.Secret{},
			wantErr: assert.Error,
			do: func() {
				utils.RefreshTestDatabase()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()

			s := &SecretPostgresStorage{conn: con}

			got, err := s.DeleteSecret(tt.args.ctx, tt.args.secret)

			if !tt.wantErr(t, err, fmt.Sprintf("DeleteSecret(%v, %v)", tt.args.ctx, tt.args.secret)) {
				assert.Equalf(t, tt.want, got, "DeleteSecret(%v, %v)", tt.args.ctx, tt.args.secret)
			}
		})
	}
}

func TestSecretPostgresStorage_EditSecret(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	uid := uuid.New()

	type args struct {
		ctx    context.Context
		secret model.Secret
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
		do      func()
	}{
		{
			name: "Secret can be edited",
			args: args{
				ctx: ctx,
				secret: model.Secret{
					ID:        1,
					UserID:    uid,
					TypeID:    1,
					Title:     "Test new",
					Content:   []byte{1, 2},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: nil,
				},
			},
			wantErr: assert.NoError,
			do: func() {
				utils.RefreshTestDatabase()

				row, _ := con.Query(
					ctx, "insert into users (id, login, password) values ($1,$2,$3)", uid, "test", "test",
				)
				row.Close()

				r2, _ := con.Query(
					ctx, "insert into secrets (user_id, title, content, type_id) values ($1,$2,$3,$4)",
					uid, "test", hex.EncodeToString([]byte{10, 20}), 1,
				)
				r2.Close()
			},
		},
		{
			name: "Error will be return when no rows is edited",
			args: args{
				ctx: ctx,
				secret: model.Secret{
					ID:        1,
					UserID:    uid,
					TypeID:    1,
					Title:     "Test",
					Content:   []byte{1, 2},
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: nil,
				},
			},
			wantErr: assert.Error,
			do: func() {
				utils.RefreshTestDatabase()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()
			s := &SecretPostgresStorage{conn: con}

			got, err := s.EditSecret(tt.args.ctx, tt.args.secret)
			if !tt.wantErr(t, err, fmt.Sprintf("EditSecret(%v, %v)", tt.args.ctx, tt.args.secret)) {
				assert.Equal(t, []byte{1, 2}, got.Content)
				assert.Equal(t, "Test new", got.Title)
				assert.NotEqual(t, time.Time{}, got.UpdatedAt)
			}
		})
	}
}

func TestSecretPostgresStorage_GetListOfSecretByType(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	uid := uuid.New()

	type args struct {
		ctx        context.Context
		secretType model.SecretType
		user       model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
		wantLen int
		do      func()
	}{
		{
			name: "List of secrets can be gotten from database by their type",
			args: args{
				ctx:        ctx,
				secretType: model.SecretType{ID: 1},
				user:       model.User{ID: &uid},
			},
			wantErr: assert.NoError,
			wantLen: 2,
			do: func() {
				utils.RefreshTestDatabase()

				row, _ := con.Query(
					ctx, "insert into users (id, login, password) values ($1,$2,$3)", uid, "test", "test",
				)
				row.Close()

				r2, _ := con.Query(
					ctx, "insert into secrets (user_id, title, content, type_id) values ($1,$2,$3,$4)",
					uid, "test", hex.EncodeToString([]byte{10, 20}), 1,
				)
				r2.Close()

				r3, _ := con.Query(
					ctx, "insert into secrets (user_id, title, content, type_id) values ($1,$2,$3,$4)",
					uid, "test2", hex.EncodeToString([]byte{1, 1}), 1,
				)
				r3.Close()
			},
		},
		{
			name: "List of secrets will return error if no rows is found",
			args: args{
				ctx:        ctx,
				secretType: model.SecretType{ID: 1},
				user:       model.User{ID: &uid},
			},
			wantLen: 0,
			wantErr: assert.NoError,
			do: func() {
				utils.RefreshTestDatabase()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()

			s := &SecretPostgresStorage{conn: con}
			got, err := s.GetListOfSecretByType(tt.args.ctx, tt.args.secretType, tt.args.user)

			tt.wantErr(t, err, fmt.Sprintf("GetListOfSecretByType(%v, %v, %v)", tt.args.ctx, tt.args.secretType, tt.args.user))

			assert.Len(t, got, tt.wantLen)
		})
	}
}
