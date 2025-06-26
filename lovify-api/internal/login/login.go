package login

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
	Username       string
}
