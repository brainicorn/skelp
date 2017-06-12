package prompter

import (
	"strconv"
)

func AsBool(answer string, err error) (bool, error) {

	if err != nil {
		return false, err
	}

	return strconv.ParseBool(answer)
}
