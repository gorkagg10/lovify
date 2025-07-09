package events

import "encoding/json"

const CreateProfile = "create-profile"

type Profile struct {
	Email string `json:"email"`
}

func NewProfile(email string) *Profile {
	return &Profile{
		Email: email,
	}
}

func (p *Profile) ToMsg() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Profile) FromMsg(data []byte) error {
	return json.Unmarshal(data, p)
}
