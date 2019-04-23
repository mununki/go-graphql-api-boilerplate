package schema

import (
	"testing"
)

func TestGetSchema(t *testing.T) {
	s := GetSchema()

	t.Log(s)
}
