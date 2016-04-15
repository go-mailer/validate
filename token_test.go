package validate

import "testing"

func TestTokenValidate(t *testing.T) {
	email := "xxx@gmail.com"
	tokenV := NewTokenValidate(NewMemoryStore(0))
	token, err := tokenV.Generate(email)
	if err != nil {
		t.Error(err)
		return
	}
	isValid, vmail, err := tokenV.Validate(token)
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
