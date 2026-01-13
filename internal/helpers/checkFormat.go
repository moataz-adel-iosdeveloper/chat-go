package helpers

import "net/mail"

func CheckEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
