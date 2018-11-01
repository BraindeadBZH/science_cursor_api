package api

// User is the model for users in database
type User struct {
	ID       int64
	Name     string
	Email    string   `sql:",unique"`
	Roles    []string `sql:",array"`
	Salt     []byte
	Password []byte
}
