package profile

import "time"

type UserProfile struct {
	email                    string
	birthday                 time.Time
	name                     string
	gender                   string
	sexualOrientation        string
	description              string
	connectedToMusicProvider bool
	musicProviderInfo        *MusicProviderData
}

func (u *UserProfile) Email() string {
	return u.email
}

func (u *UserProfile) Birthday() time.Time {
	return u.birthday
}

func (u *UserProfile) Gender() string {
	return u.gender
}

func (u *UserProfile) Name() string {
	return u.name
}

func (u *UserProfile) SexualOrientation() string {
	return u.sexualOrientation
}

func (u *UserProfile) Description() string {
	return u.description
}

func (u *UserProfile) MusicProviderInfo() *MusicProviderData {
	return u.musicProviderInfo
}

func (u *UserProfile) ConnectedToMusicProvider() bool {
	return u.connectedToMusicProvider
}

func NewUserProfile(
	email string,
	birthday time.Time,
	name string,
	gender string,
	sexualOrientation string,
	description string,
) *UserProfile {
	return &UserProfile{
		email:                    email,
		birthday:                 birthday,
		name:                     name,
		gender:                   gender,
		sexualOrientation:        sexualOrientation,
		description:              description,
		connectedToMusicProvider: false,
	}
}
