package tests

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RequireValidationErrors(t *testing.T, expected map[string]string, err error) {
	var errors validator.ValidationErrors
	require.ErrorAs(t, err, &errors)
	assert.Len(t, errors, len(expected))

	actualErrors := make(map[string]string)
	for _, fieldError := range errors {
		actualErrors[fieldError.Field()] = fieldError.Tag()
	}
	for field, tag := range expected {
		if actualTag, ok := actualErrors[field]; ok {
			assert.Equal(t, tag, actualTag)
		} else {
			assert.Fail(t, "unexpected validation error", "field: %s, tag: %s", field, tag)
		}
	}
}
