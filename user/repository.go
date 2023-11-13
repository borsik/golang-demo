package user

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Repository interface {
	Insert(input InputUser) (uuid.UUID, error)
	Select(name string, country string, offset int, limit int) ([]User, int64, error)
	SelectById(id uuid.UUID) (User, error)
	Update(id uuid.UUID, input InputUser) error
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

var psql sq.StatementBuilderType

func init() {
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func (r *repository) Insert(input InputUser) (uuid.UUID, error) {
	var id uuid.UUID
	query :=
		psql.Insert("users").SetMap(map[string]interface{}{
			"first_name": input.FirstName,
			"last_name":  input.LastName,
			"nickname":   input.Nickname,
			"password":   input.Password,
			"email":      input.Email,
			"country":    input.Country,
		}).Suffix("RETURNING id")
	err := query.RunWith(r.db).QueryRow().Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (r *repository) Select(name string, country string, offset int, limit int) ([]User, int64, error) {
	var users []User
	var totalCount int64

	builder := func(sel sq.SelectBuilder) sq.SelectBuilder {
		query := sel.From("users")
		if name != "" {
			query = query.Where("first_name ILIKE ? OR last_name ILIKE ?",
				fmt.Sprint("%", name, "%"),
				fmt.Sprint("%", name, "%"))
		}
		if country != "" {
			query = query.Where("country = ?", strings.ToUpper(country))
		}
		return query
	}

	err := builder(psql.Select("count(1) AS total")).RunWith(r.db).QueryRow().Scan(&totalCount)
	if err != nil {
		return users, totalCount, err
	}
	rows, err := builder(psql.Select("id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at")).
		Offset(uint64(offset)).Limit(uint64(limit)).RunWith(r.db).Query()
	if err != nil {
		return users, totalCount, err
	}
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Password, &u.Email, &u.Country, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return users, totalCount, err
		}
		users = append(users, u)
	}

	if err != nil {
		return users, totalCount, err
	}
	return users, totalCount, nil
}

func (r *repository) SelectById(id uuid.UUID) (User, error) {
	var u User
	query :=
		psql.Select("id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at").
			From("users").Where(sq.Eq{"id": id})
	err := query.RunWith(r.db).QueryRow().Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Password, &u.Email, &u.Country, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (r *repository) Update(id uuid.UUID, input InputUser) error {
	query := psql.Update("users").SetMap(map[string]interface{}{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"nickname":   input.Nickname,
		"password":   input.Password,
		"email":      input.Email,
		"country":    input.Country,
		"updated_at": time.Now(),
	}).Where("id = ?", id)
	_, err := query.RunWith(r.db).Exec()
	return err
}

func (r *repository) Delete(id uuid.UUID) error {
	query := psql.Delete("users").Where("id = ?", id)
	_, err := query.RunWith(r.db).Exec()
	return err
}
