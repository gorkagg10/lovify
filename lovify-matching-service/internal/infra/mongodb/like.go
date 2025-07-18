package mongodb

import "time"

type Like struct {
	FromUserId string    `bson:"from_user_id"`
	ToUserId   string    `bson:"to_user_id"`
	Type       string    `bson:"type"`
	CreatedAt  time.Time `bson:"created_at"`
}
