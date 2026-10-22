package implementation

import (
	"errors"

	"security/internal/repository"
	"security/internal/service"
)

func translate(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return service.ErrNotFound
	}
	return err
}
