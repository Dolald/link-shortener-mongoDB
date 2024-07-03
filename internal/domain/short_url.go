package pkg

import "time"

type ShortURL struct {
	Id        string     `bson:"_id"`
	URL       string     `bson:"url"`
	ExpiredAt *time.Time `bson:"expireAt,omitempty"`
}
