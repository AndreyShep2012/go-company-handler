package repositories

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type ErrNotFound struct{}

func (ErrNotFound) Error() string {
	return "not found"
}

type ErrDuplicatedKey struct{}

func (ErrDuplicatedKey) Error() string {
	return "duplicated key"
}

func handleError(err error) error {
	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicatedKey{}
	}

	switch err {
	case mongo.ErrNoDocuments:
		return errors.Join(ErrNotFound{}, err)
	default:
		return err
	}
}
