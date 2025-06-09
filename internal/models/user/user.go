package user

import (
	"time"
)

type User struct {
	Id        int64
	Username  string
	CreatedAt time.Time
}
