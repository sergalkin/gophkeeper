package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	"github.com/sergalkin/gophkeeper/pkg/utils"
)

func TestNewPostgresSecretTypeStorage(t *testing.T) {
	tests := []struct {
		name string
		want *SecretTypePostgresStorage
	}{
		{
			name: "New Postgres Secret Type Storage could be created",
			want: &SecretTypePostgresStorage{&pgx.Conn{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPostgresSecretTypeStorage(&pgx.Conn{})
			assert.Equalf(t, tt.want, got, "NewPostgresSecretTypeStorage() = %v, want %v", got, tt.want)
		})
	}
}

func Test_secretTypePostgresStorage_GetSecretTypes(t *testing.T) {
	ctx := context.Background()

	con := utils.CreatePostgresTestConn()
	defer con.Close(ctx)

	tests := []struct {
		name  string
		empty bool
		do    func()
	}{
		{
			name:  "Secret types list can be gotten from database",
			empty: false,
			do: func() {
				utils.RefreshTestDatabase()
			},
		},
		{
			name:  "Secret type list will be empty on retrieval",
			empty: true,
			do: func() {
				r, _ := con.Query(ctx, "delete from secret_types")
				r.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.do()

			s := &SecretTypePostgresStorage{
				conn: con,
			}

			got, err := s.GetSecretTypes(ctx)
			if err == nil && tt.empty == false {
				assert.NotEmpty(t, got)
				return
			}

			assert.Empty(t, got)
		})
	}
}
