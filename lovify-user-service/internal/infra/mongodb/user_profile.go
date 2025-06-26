package mongodb

type UserProfile struct {
	ID                     string `bson:"_id,omitempty"`
	Email                  string `bson:"email"`
	Name                   string `bson:"name"`
	Birthday               string `bson:"birthday"`
	Gender                 string `bson:"gender"`
	SexualOrientation      string `bson:"sexual_orientation"`
	Description            string `bson:"description"`
	MusicProviderConnected bool   `bson:"music_provider_connected"`
}

func NewUserProfile(
	id string,
	email string,
	name string,
	birthday string,
	gender string,
	sexualOrientation string,
	description string,
	musicProviderConnected bool,
) *UserProfile {
	return &UserProfile{
		ID:                     id,
		Email:                  email,
		Name:                   name,
		Birthday:               birthday,
		Gender:                 gender,
		SexualOrientation:      sexualOrientation,
		Description:            description,
		MusicProviderConnected: musicProviderConnected,
	}
}
