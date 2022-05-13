package models

import "time"

type Stock struct {
	Id        int32
	Name      string
	UserId    int32
	Amount    int32
	CreatedAt time.Time
}
