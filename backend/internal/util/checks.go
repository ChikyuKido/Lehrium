package util

import (
    "regexp"
)

func CheckForEmailDomain(email string) (bool, error) {
	email_match, err := regexp.MatchString(`@spengergasse\.at$`, email)
	if !email_match {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
