package hasher

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}
