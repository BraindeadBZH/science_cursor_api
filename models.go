package api

// User is the model for users in database
type User struct {
	ID       int64
	Name     string
	Email    string
	Roles    []string `sql:",array"`
	Salt     []byte
	Password []byte
}
