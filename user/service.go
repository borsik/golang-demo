package user

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	Store(input InputUser) (uuid.UUID, error)
	Get(name string, country string, page int, pageSize int) ([]User, int64, error)
	GetById(id uuid.UUID) (User, error)
	Update(id uuid.UUID, input InputUser) error
	Delete(id uuid.UUID) error
}

type service struct {
	repository Repository
	amqp       MQ
	log        log.FieldLogger
}

func NewService(repository Repository, amqp MQ, logger log.FieldLogger) *service {
	return &service{repository, amqp, logger}
}

func (s *service) Store(input InputUser) (uuid.UUID, error) {
	id, err := s.repository.Insert(
		User{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Nickname:  input.Nickname,
			Password:  input.Password,
			Email:     input.Email,
			Country:   input.Country})
	if err != nil {
		return id, err
	}
	s.amqp.PublishMessage("create_user", id.String())
	s.log.Debugln("created user", id)
	return id, nil
}

func (s *service) Get(name string, country string, page int, pageSize int) ([]User, int64, error) {
	offset := (page - 1) * pageSize
	limit := pageSize
	users, totalCount, err := s.repository.Select(name, country, offset, limit)
	return users, totalCount, err
}

func (s *service) GetById(id uuid.UUID) (User, error) {
	user, err := s.repository.SelectById(id)
	return user, err
}

func (s *service) Update(id uuid.UUID, input InputUser) error {
	err := s.repository.Update(id, input)
	s.amqp.PublishMessage("update_user", id.String())
	s.log.Debugln("updated user", id)
	return err
}

func (s *service) Delete(id uuid.UUID) error {
	err := s.repository.Delete(id)
	s.amqp.PublishMessage("delete_user", id.String())
	s.log.Debugln("deleted user", id)
	return err
}
