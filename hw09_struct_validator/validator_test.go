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
		expectedErr error
	}{
		{
			in: Response{
				Code: 200,
				Body: "lalala",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 999,
				Body: "lalala",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrConstraintCheckFailed,
				},
			},
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in: App{Version: "123456"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrConstraintCheckFailed,
				},
			},
		},
		{
			in: App{Version: "1234"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrConstraintCheckFailed,
				},
			},
		},
		{
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "Name",
				Age:    21,
				Email:  "test@mil.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "012345678901234567890123456789012345",
				Name:   "Name",
				Age:    21,
				Email:  "test.email@mail.com@email",
				Role:   "admin",
				Phones: []string{"12345678901"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrConstraintCheckFailed,
				},
			},
		},
		{
			in: User{
				ID:     "___",
				Name:   "Name",
				Age:    15,
				Email:  "test.email@mail.com@email",
				Role:   "noinlist",
				Phones: []string{"12345678901", "shortphone"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrConstraintCheckFailed,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrConstraintCheckFailed,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrConstraintCheckFailed,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrConstraintCheckFailed,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrConstraintCheckFailed,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			resError := Validate(tt.in)
			var validationErrors ValidationErrors
			if errors.As(resError, &validationErrors) {
				require.Equal(t, tt.expectedErr, validationErrors)
				return
			}
			require.Equal(t, resError, tt.expectedErr)
		})
	}
}
