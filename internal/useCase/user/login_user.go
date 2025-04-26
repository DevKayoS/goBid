package user

import (
	"context"

	"github.com/DevKayoS/goBid/internal/validator"
)

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.Matches(req.Email, validator.EmailRx), "email", "must be a valid email")
	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")

	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")

	return eval
}
