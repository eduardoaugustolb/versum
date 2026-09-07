package domain

type User struct {
	id                        string
	emailCiphertext           []byte
	emailLookupHMAC           []byte
	emailEncryptionKeyVersion int
	emailLookupKeyVersion     int
}

func NewUser(id string, emailCiphertext []byte, emailLookupHMAC []byte, emailEncryptionKeyVersion int, emailLookupKeyVersion int) (*User, error) {

	user := &User{
		id:                        id,
		emailCiphertext:           cloneBytes(emailCiphertext),
		emailLookupHMAC:           cloneBytes(emailLookupHMAC),
		emailEncryptionKeyVersion: emailEncryptionKeyVersion,
		emailLookupKeyVersion:     emailLookupKeyVersion,
	}

	if err := validateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func RehydrateUser(id string, emailCiphertext []byte, emailLookupHMAC []byte, emailEncryptionKeyVersion int, emailLookupKeyVersion int) (*User, error) {
	user := &User{
		id:                        id,
		emailCiphertext:           cloneBytes(emailCiphertext),
		emailLookupHMAC:           cloneBytes(emailLookupHMAC),
		emailEncryptionKeyVersion: emailEncryptionKeyVersion,
		emailLookupKeyVersion:     emailLookupKeyVersion,
	}

	if err := validateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func validateUser(user *User) error {
	if user.id == "" {
		return ErrInvalidUserID
	}

	if len(user.emailCiphertext) == 0 {
		return ErrInvalidEmailCiphertext
	}

	if len(user.emailLookupHMAC) == 0 {
		return ErrInvalidEmailLookupHMAC
	}

	if user.emailEncryptionKeyVersion <= 0 {
		return ErrInvalidEmailEncryptionKeyVersion
	}

	if user.emailLookupKeyVersion <= 0 {
		return ErrInvalidEmailLookupKeyVersion
	}

	return nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) EmailCiphertext() []byte {
	return cloneBytes(u.emailCiphertext)
}

func (u *User) EmailLookupHMAC() []byte {
	return cloneBytes(u.emailLookupHMAC)
}

func (u *User) EmailEncryptionKeyVersion() int {
	return u.emailEncryptionKeyVersion
}

func (u *User) EmailLookupKeyVersion() int {
	return u.emailLookupKeyVersion
}
