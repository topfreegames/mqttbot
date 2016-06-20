package modules

import (
	"strings"
	"testing"
)

func TestGenHash(t *testing.T) {
	generated := GenHash("password")
	if !strings.HasPrefix(generated, "PBKDF2$sha256$901$") {
		t.Fail()
	}
}
