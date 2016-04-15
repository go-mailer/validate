package validate

import "testing"

func TestCodeValidate(t *testing.T) {
	email := "xxx@gmail.com"
	codeV := NewCodeValidate(NewMemoryStore(0))
	code, err := codeV.Generate(email)
	if err != nil {
		t.Error(err)
		return
	}
	isValid, err := codeV.Validate(email, code)
	if err != nil {
		t.Error(err)
		return
	}
	if !isValid {
		t.Error("Validate error")
		return
	}
}
