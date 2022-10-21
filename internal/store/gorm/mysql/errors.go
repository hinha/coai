package mysql

import (
	"github.com/pkg/errors"
)

var (
	DataAlreadyExists = errors.New("data already exists")
	RecordNotFound    = errors.New("data not found")
)
