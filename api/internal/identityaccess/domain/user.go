package domain

type User struct {
	id    string
	email Email
}

func NewUser(id string, email Email) (*User, error) {
	user := &User{id: id, email: email}

	if err := validateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func RehydrateUser(id string, email Email) (*User, error) {
	user := &User{id: id, email: email}

	if err := validateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func validateUser(user *User) error {
	if user.id == "" {
		return ErrInvalidUserID
	}

	return nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() Email {
	return u.email
}
