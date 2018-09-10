package prompter

import (
	"fmt"
	"testing"
)

var asBoolTests = []struct {
	answer        string
	expectedBool  bool
	err           error
	expectedError error
}{
	{
		"true",
		true,
		nil,
		nil,
	},
	{
		"t",
		true,
		nil,
		nil,
	},
	{
		"1",
		true,
		nil,
		nil,
	},
	{
		"true",
		false,
		fmt.Errorf("forced error"),
		fmt.Errorf("forced error"),
	},
	{
		"blah",
		false,
		nil,
		fmt.Errorf(`strconv.ParseBool: parsing "blah": invalid syntax`),
	},
	{
		"false",
		false,
		nil,
		nil,
	},
	{
		"f",
		false,
		nil,
		nil,
	},
	{
		"0",
		false,
		nil,
		nil,
	},
	{
		"false",
		false,
		fmt.Errorf("forced error"),
		fmt.Errorf("forced error"),
	},
	{
		"blah",
		false,
		nil,
		fmt.Errorf(`strconv.ParseBool: parsing "blah": invalid syntax`),
	},
}

func TestAsBool(t *testing.T) {
	for _, abt := range asBoolTests {
		errString, expErrorString := "nil", "nil"
		b, err := AsBool(abt.answer, abt.err)

		if err != nil {
			errString = err.Error()
		}

		if abt.expectedError != nil {
			expErrorString = abt.expectedError.Error()
		}

		if errString != expErrorString {
			t.Errorf("ABT errors don't match have (%s) want (%s)", errString, expErrorString)
		}

		if b != abt.expectedBool {
			t.Errorf("ABT bools don't match have (%t) want (%t)", b, abt.expectedBool)
		}
	}
}
