package user

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func TestFindUser(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	id := uuid.New()
	users := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country"}).
		AddRow(id, "first_name", "last_name", "nickname", "passwd", "asd@mail.ru", "kz")

	expectedSQL := "SELECT (.+) FROM \"users\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	_, res := repo.SelectById(id)

	assert.Nil(t, res)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestAddUser(t *testing.T) {
	sqlDB, db, _ := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	user := User{}
	addedUser, _ := repo.Insert(user)
	assert.Equal(t, addedUser.ID, user.ID)
}

func TestDeleteUser(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	id := uuid.New()
	delSQL := "DELETE FROM \"users\" WHERE \"users\".\"id\" = .+"
	mock.ExpectBegin()
	mock.ExpectExec(delSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := repo.Delete(id)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewRepository(db)

	updUserSQL := "UPDATE \"users\" SET .+"
	mock.ExpectBegin()
	mock.ExpectExec(updUserSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	id := uuid.New()
	err := repo.Update(id, InputUser{FirstName: "name"})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
