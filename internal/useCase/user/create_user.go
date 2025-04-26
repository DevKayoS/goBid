package user

import (
	"context"

	"github.com/DevKayoS/goBid/internal/validator"
)

type CreateUserRequest struct {
	UserName     string `json:"user_name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Bio          string `json:"bio"`
}

func (req CreateUserRequest) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	// validate stuff
	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")
}
