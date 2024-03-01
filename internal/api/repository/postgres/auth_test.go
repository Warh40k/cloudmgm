package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Warh40k/cloud-manager/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthPostgres_SignUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	r := NewAuthPostgres(dbx)
	id := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		input   domain.User
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "ValidCredentials",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs(id, "Test", "test", "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO").WillReturnRows(rows)
			},
			input: domain.User{
				Id:           id,
				Name:         "Test",
				Username:     "test",
				PasswordHash: "$2a$10$7wyI.VRyw8GRxUBp9Gi3b.S7EH6u45HtKeG3GklSkSLtpoceXYAlO",
			},
			want: id,
		},
		{
			name: "EmptyPassword",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", usersTable)).
					WithArgs(id, "Test", "test", "").WillReturnRows(rows)
			},
			input: domain.User{
				Id:       id,
				Name:     "Test",
				Username: "test",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.SignUp(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
