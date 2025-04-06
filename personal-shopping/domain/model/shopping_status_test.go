package model_test

import (
	"errors"
	"share-basket/personal-shopping/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoppingStatus(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    model.ShoppingStatus
		expectedErr error
	}{
		"正常系": {
			input:       "un_purchased",
			expected:    model.UnPurchased,
			expectedErr: nil,
		},
		"不正なステータスが渡された場合": {
			input:       "invalid_status",
			expected:    "",
			expectedErr: errors.New("invalid shopping status"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			status, err := model.NewShoppingStatus(tt.input)

			assert.Equal(t, tt.expected, status)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
