package mongodb

import "time"

type Match struct {
	ID                  string    `json:"id" bson:"id"`
	User1ID             string    `bson:"user_1_id"`
	User2ID             string    `bson:"user_2_id"`
	MatchedAt           time.Time `bson:"matched_at"`
	ConversationStarted bool      `bson:"conversation_started"`
}
