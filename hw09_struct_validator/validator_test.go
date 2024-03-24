package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []error
	}{
		{
			User{
				ID:     "100",
				Name:   "John",
				Age:    16,
				Email:  "test@test.ru",
				Role:   "stuff",
				Phones: []string{"5466", "65464584836"},
				meta:   nil,
			},
			[]error{
				ErrLenViolation,
				ErrMinValueViolation,
				ErrLenViolation,
			},
		},
		{
			App{
				Version: "1234",
			},
			[]error{
				ErrLenViolation,
			},
		},
		{
			Token{},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "test",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var valErrs ValidationErrors

			if errors.As(err, &valErrs) {
				for i, err := range tt.expectedErr {
					require.ErrorIs(t, valErrs[i].Err, err, "Validation error should be like expected")
				}
			}
		})
	}

	t.Run("Should handle non struct value", func(t *testing.T) {
		err := Validate(99999)
		require.ErrorIs(t, err, ErrExpectedStruct)
	})

	t.Run("Should handle nil value", func(t *testing.T) {
		err := Validate(nil)
		require.ErrorIs(t, err, ErrExpectedStruct)
	})
}
