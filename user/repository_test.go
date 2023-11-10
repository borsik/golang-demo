package user

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func DbMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, mock
}

func TestFindUserById(t *testing.T) {
	db, mock := DbMock(t)
	defer db.Close()
	repo := NewRepository(db)

	id := uuid.New()
	users := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
		AddRow(id, "firstname", "lastname", "nickname", "passwd", "example@mail.com", "xx", time.Now(), time.Now())

	expectedSQL := "SELECT (.+) FROM users WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)
	_, err := repo.SelectById(id)

	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestFindUser(t *testing.T) {
	db, mock := DbMock(t)
	defer db.Close()
	repo := NewRepository(db)

	expectedCount := "SELECT count(.+) AS total FROM users WHERE first_name ILIKE (.+) OR last_name ILIKE (.+)"
	expectedSelect := "SELECT (.+) FROM users WHERE first_name ILIKE (.+) OR last_name ILIKE (.+) LIMIT (.+) OFFSET (.+)"

	mock.ExpectQuery(expectedCount).WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(0))

	usersRow := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
		AddRow(uuid.New(), "firstname", "lastname", "nickname", "passwd", "example@mail.com", "xx", time.Now(), time.Now())
	mock.ExpectQuery(expectedSelect).WillReturnRows(usersRow)

	_, _, err := repo.Select("name", "", 0, 1)

	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)
}

func TestAddUser(t *testing.T) {
	db, mock := DbMock(t)
	defer db.Close()
	repo := NewRepository(db)

	expectedQuery := "INSERT INTO users (.+) VALUES (.+) RETURNING id"
	mock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

	_, err := repo.Insert(InputUser{FirstName: "first", LastName: "last", Nickname: "nick"})
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock := DbMock(t)
	defer db.Close()
	repo := NewRepository(db)

	expectedSQL := "DELETE FROM users WHERE id = (.+)"
	mock.ExpectExec(expectedSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.Delete(uuid.New())
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	db, mock := DbMock(t)
	defer db.Close()
	repo := NewRepository(db)

	expectedSQL := "UPDATE users SET (.+) WHERE id = (.+)"
	mock.ExpectExec(expectedSQL).WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.Update(uuid.New(), InputUser{FirstName: "name"})
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
