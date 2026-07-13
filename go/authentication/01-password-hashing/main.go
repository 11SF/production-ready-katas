package main

import (
	"errors"
	"go/authentication/01-password-hashing/password-hasher"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	pwd := "my-secret-password"

	pwdHasher := password.NewPasswordHasher()

	hashPassword, err := pwdHasher.HashPassword(pwd)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			println("Password is too long (> 72 bytes).")
		}
		panic(err)
	}

	println("Hash:", hashPassword)

	isValid := pwdHasher.VerifyPassword(hashPassword, pwd)
	println("Is valid password:", isValid)
}
