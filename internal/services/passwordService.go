package services

import "golang.org/x/crypto/bcrypt"

func HashPassword(password, pepper string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hashedPassword, password, pepper string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+pepper))
	if err != nil {
		return err
	}
	return nil
}
