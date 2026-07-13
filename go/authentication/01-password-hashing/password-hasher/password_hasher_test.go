package password_test

import (
	"go/authentication/01-password-hashing/password-hasher"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHasher(t *testing.T) {
	t.Run("cost >= 12", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		assert.True(t, pwdHasher.Cost >= 12)
	})

	t.Run("HashPassword('correct-password') คืน hash ที่ VerifyPassword ตรวจสอบผ่าน", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		pwd := "correct-password"
		hash, err := pwdHasher.HashPassword(pwd)
		assert.NoError(t, err)
		assert.True(t, pwdHasher.VerifyPassword(hash, pwd))
	})

	t.Run("VerifyPassword(hash, 'wrong-password') คืน false — ไม่ panic ไม่ return error", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		pwd := "correct-password"
		hash, err := pwdHasher.HashPassword(pwd)
		assert.NoError(t, err)
		assert.False(t, pwdHasher.VerifyPassword(hash, "wrong-password"))
	})

	t.Run("HashPassword สองครั้งด้วย password เดียวกัน คืน hash ที่ต่างกัน (salt ต่างกัน)", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		pwd := "correct-password"
		hash1, err := pwdHasher.HashPassword(pwd)
		assert.NoError(t, err)
		hash2, err := pwdHasher.HashPassword(pwd)
		assert.NoError(t, err)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("HashPassword ด้วย password ที่ยาวกว่า 72 bytes คืน error ที่ชัดเจน", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		pwd := "this-password-is-way-too-long-and-should-trigger-an-error-because-it-exceeds-the-maximum-length-of-72-bytes"
		_, err := pwdHasher.HashPassword(pwd)
		assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})

	t.Run("Hash ที่ได้มี prefix $2a$ และ cost >= 12 เมื่อดูด้วยตา", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		pwd := "Secret"
		hash, err := pwdHasher.HashPassword(pwd)

		assert.True(t, strings.HasPrefix(hash, "$2a$"))
		assert.True(t, strings.Contains(hash, "$12$"))
		assert.NoError(t, err)
	})

	t.Run("VerifyPassword ด้วย hash ที่ malformed คืน false — ไม่ crash", func(t *testing.T) {
		pwdHasher := password.NewPasswordHasher()
		assert.False(t, pwdHasher.VerifyPassword("malformed-hash", "any-password"))
	})
}
