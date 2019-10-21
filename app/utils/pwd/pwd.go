package pwd

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(src string) (pwd string, err error) {
	raw, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	return string(raw), err
}

func ComparePassword(hashedPwd string, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	return err == nil
}
