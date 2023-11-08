package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Service interface {
	Store(input InputUser) (User, error)
	Get(name string, country string, page int, pageSize int) ([]User, int64)
	GetByID(id string) (User, error)
	Update(id string, input InputUser) error
	Delete(id string) error
}

type service struct {
	repository Repository
	amqp       MQ
	log        log.FieldLogger
}

func NewService(repository Repository, amqp MQ, logger log.FieldLogger) *service {
	return &service{repository, amqp, logger}
}

func (s *service) Store(input InputUser) (User, error) {
	user := User{ID: uuid.New().String(), FirstName: input.FirstName, LastName: input.LastName, Nickname: input.Nickname,
		Password: input.Password, Email: input.Email, Country: strings.ToUpper(input.Country)}
	newUser, err := s.repository.Insert(user)
	if err != nil {
		return newUser, err
	}
	s.amqp.PublishMessage("create_user", user.ID)
	s.log.Debugln("created user", user.ID)
	return user, nil
}

func (s *service) Get(name string, country string, page int, pageSize int) ([]User, int64) {
	offset := (page - 1) * pageSize
	limit := pageSize
	users, totalCount := s.repository.Select(name, country, offset, limit)
	return users, totalCount
}

func (s *service) GetByID(id string) (User, error) {
	user, err := s.repository.SelectByID(id)
	return user, err
}

func (s *service) Update(id string, input InputUser) error {
	err := s.repository.Update(id, input)
	s.amqp.PublishMessage("update_user", id)
	s.log.Debugln("updated user", id)
	return err
}

func (s *service) Delete(id string) error {
	err := s.repository.Delete(id)
	s.amqp.PublishMessage("delete_user", id)
	s.log.Debugln("deleted user", id)
	return err
}
