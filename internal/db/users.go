package db

type User struct {
	Email string `bson:"email"`

	SessionExpires string `bson:"sessionExpires"`
}

func (c *Connection) GetUserBySessionToken(sessionToken string) (*User, error) {

	// filter := bson.M{"sessionToken": sessionToken}
	return nil, nil
}
