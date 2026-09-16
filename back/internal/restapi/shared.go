package restapi

import (
	"errors"
)

const (
	LocalsUserIDKey = "userID"
)

var (
	ErrGetUserID = errors.New("get user id")
)
