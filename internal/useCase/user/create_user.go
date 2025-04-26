package user

import (
	"context"

	"github.com/DevKayoS/goBid/internal/validator"
)

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	// validate stuff
	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")

	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRx), "email", "must be a valid email")

	eval.CheckField(validator.NotBlank(req.Bio), "bio", "this field cannot be empty")

	eval.CheckField(
		validator.MinChars(req.Bio, 10) &&
			validator.MaxChars(req.Bio, 255), "bio", "this field must have a length between 10 and 255",
	)

	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.Password, 8), "password", "Password must be at least 8 characters long.")

	return eval
}
