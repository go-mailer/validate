package validate

import (
	"testing"
)

func TestManagerValidate(t *testing.T) {
	email := "xxx@gmail.com"
	manager := NewManager(NewMemoryStore(0))
	token, err := manager.GenerateToken(email)
	if err != nil {
		t.Error(err)
		return
	}
	isValid, vmail, err := manager.Validate(token)
	if err != nil {
		t.Error(err)
		return
	}
	if !isValid {
		t.Error("Validate error")
		return
	}
	if vmail != email {
		t.Error("Get email error")
	}
}
