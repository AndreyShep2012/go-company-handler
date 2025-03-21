package services

import (
	"errors"

	"github.com/AndreyShep2012/go-company-handler/internal/app/v1/repositories"
)

type ErrNotFound struct{}

func (ErrNotFound) Error() string {
	return "not found"
}

type ErrDb struct{}

func (ErrDb) Error() string {
	return "db error"
}

type ErrDbDuplicatedKey struct{}

func (ErrDbDuplicatedKey) Error() string {
	return "db duplicated key"
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.As(err, &repositories.ErrNotFound{}):
		return errors.Join(ErrNotFound{}, err)
	case errors.As(err, &repositories.ErrDuplicatedKey{}):
		return errors.Join(ErrDbDuplicatedKey{}, err)
	default:
		return errors.Join(ErrDb{}, err)
	}
}
