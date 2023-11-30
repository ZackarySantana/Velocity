package db

type User struct {
	Email string `bson:"email"`

	SessionExpires string `bson:"sessionExpires"`
}

func (u *User) CheckPassword(password string) bool {
	return true
}

func (c *Connection) GetUserBySessionToken(sessionToken string) (*User, error) {

	// filter := bson.M{"sessionToken": sessionToken}
	return nil, nil
}

func (c *Connection) GetUserByEmail(username string) (*User, error) {

	// filter := bson.M{"username": username}
	return nil, nil
}
