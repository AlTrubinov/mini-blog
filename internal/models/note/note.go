package note

import "time"

type Note struct {
	Id        int64
	UserId    int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
