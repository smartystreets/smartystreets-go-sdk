package contracts

import "errors"

var (
	ErrNoAuthorization = errors.New("authorization has not been configured")
)