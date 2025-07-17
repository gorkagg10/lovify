package events

import "encoding/json"

const CreateProfile = "create-profile"

type Profile struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
}

func NewProfile(email, userID string) *Profile {
	return &Profile{
		Email:  email,
		UserID: userID,
	}
}

func (p *Profile) ToMsg() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Profile) FromMsg(data []byte) error {
	return json.Unmarshal(data, p)
}
