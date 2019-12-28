package pwd

import "testing"

func TestPassword(t *testing.T) {
	gotPwd, err := GeneratePassword("hello pwd")
	if err != nil {
		t.Fatal(err)
	}

	ok := ComparePassword(gotPwd, "hello pwd")
	if !ok {
		t.Errorf("pwd faild")
	}
}
