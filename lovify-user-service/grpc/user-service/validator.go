package service

import "github.com/gorkagg10/lovify-user-service/errors"

func (u *CreateUserRequest) Validate() error {
	switch {
	case u == nil:
		return errors.ErrInvalidCreateUserRequest
	case u.GetEmail() == "":
		return errors.ErrInvalidEmail
	case u.GetBirthday() == nil:
		return errors.ErrInvalidBirthday
	case u.GetGender() == Gender_UNKNOWN_GENDER:
		return errors.ErrInvalidCreateUserRequest
	case u.GetSexualOrientation() == SexualOrientation_UNKNOWN_SEXUAL_ORIENTATION:
		return errors.ErrInvalidCreateUserRequest
	case u.GetName() == "":
		return errors.ErrInvalidName
	default:
		return nil
	}
}
