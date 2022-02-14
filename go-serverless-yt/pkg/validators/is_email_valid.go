package validators

import "regexp"

func IsEmailValid(email string) bool {
	// 정규표현식을 컴파일해서 뭐 한다는디 잘 모르겠네
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{1,64}@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// 아마 컴파일된 정규표현식과 email이 일치해야만 Validation을 통과 하는듯
	if len(email) < 3 || len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}

	return true
}
