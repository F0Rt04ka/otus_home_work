package hw09structvalidator

import (
	"encoding/json"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
			in: User{
				ID:     "123",
				Age:    17,
				Email:  "test@mail.ru",
				Role:   "test",
				Phones: []string{"123"},
			},
			expectedErr: fmt.Errorf("ID: value length must be equal to 36, " +
				"Age: value must be greater than 18, " +
				"Role: value must be one of admin,stuff, " +
				"Phones[0]: value length must be equal to 11"),
		},
		{
			in: User{
				ID:    "123456789012345678901234567890123456",
				Age:   20,
				Email: "test@mail.ru",
				Role:  "admin",
				Phones: []string{
					"12345678901",
				},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "123456",
			},
			expectedErr: fmt.Errorf("Version: value length must be equal to 5"),
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "body",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
