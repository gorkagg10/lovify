package postgres

type User struct {
	ID                 int64
	Email              string
	Password           string
	IsProfileConnected bool
}

type Token struct {
	ID             int64
	Token          string
	TokenType      string
	ExpirationDate string
	Email          string
}
