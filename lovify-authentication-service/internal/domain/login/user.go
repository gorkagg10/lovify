package login

type User struct {
	id                 string
	email              string
	hashedPassword     string
	sessionToken       *Token
	csrfToken          *Token
	isProfileConnected bool
	profileID          string
}

func NewUser(email string, hashedPassword string, isProfileConnected bool, profileID string) *User {
	return &User{
		email:              email,
		hashedPassword:     hashedPassword,
		isProfileConnected: isProfileConnected,
		profileID:          profileID,
	}
}

func (u *User) Email() string {
	return u.email
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

func (u *User) SessionToken() *Token {
	return u.sessionToken
}

func (u *User) CSRFToken() *Token {
	return u.csrfToken
}

func (u *User) setSessionToken(sessionToken *Token) {
	u.sessionToken = sessionToken
}

func (u *User) setCSRFToken(csrfToken *Token) {
	u.csrfToken = csrfToken
}

func (u *User) IsProfileConnected() bool {
	return u.isProfileConnected
}

func (u *User) ProfileID() string {
	return u.profileID
}
