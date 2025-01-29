package utils

import "regexp"

var (
	reEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rePhone = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func ValidateEmail(email string) bool {
	return reEmail.MatchString(email)
}

func ValidatePhone(phone string) bool {
	return rePhone.MatchString(phone)
}
