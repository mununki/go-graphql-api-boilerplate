package schema

import (
	"testing"
)

func TestGetSchema(t *testing.T) {
	s := NewSchema()

	t.Log(*s)
}
