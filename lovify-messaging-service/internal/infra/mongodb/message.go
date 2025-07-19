package mongodb

import "time"

type Message struct {
	MatchID    string    `bson:"match_id"`
	FromUserID string    `bson:"from_user_id"`
	ToUserID   string    `bson:"to_user_id"`
	Content    string    `bson:"content"`
	SentAt     time.Time `bson:"sent_at"`
	Read       bool      `bson:"read"`
}
