package password

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct {
	Cost int
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{Cost: 12}
}

func (ph *PasswordHasher) HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), ph.Cost)
	return string(hash), err
}

func (ph *PasswordHasher) VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
