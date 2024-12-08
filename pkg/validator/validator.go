package validator

import (
	"context"

	"github.com/go-playground/validator"
)

func Validate(ctx context.Context, dto any) error {
	validate := validator.New()

	err := validate.StructCtx(ctx, dto)

	if err == nil {
		return nil
	}

	return err
}
