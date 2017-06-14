package prompter

import (
	"fmt"
	"testing"
)

func TestStringNotBlank(t *testing.T) {
	if StringNotBlank("woohoo") != nil {
		t.Errorf("should not have gotten a string blank error")
	}
}

func TestStringNotBlankBlank(t *testing.T) {
	if StringNotBlank("   ") == nil {
		t.Errorf("should have gotten a string blank error")
	}
}

var ynTests = []struct {
	in       string
	expected error
}{
	{
		"Y",
		nil,
	},
	{
		"y",
		nil,
	},
	{
		"yes",
		nil,
	},
	{
		"YES",
		nil,
	},
	{
		"yEs",
		nil,
	},
	{
		"T",
		nil,
	},
	{
		"t",
		nil,
	},
	{
		"True",
		nil,
	},
	{
		"true",
		nil,
	},
	{
		"tRuE",
		nil,
	},
	{
		"N",
		nil,
	},
	{
		"n",
		nil,
	},
	{
		"no",
		nil,
	},
	{
		"NO",
		nil,
	},
	{
		"nO",
		nil,
	},
	{
		"F",
		nil,
	},
	{
		"f",
		nil,
	},
	{
		"True",
		nil,
	},
	{
		"false",
		nil,
	},
	{
		"fAlsE",
		nil,
	},
	{
		"1",
		nil,
	},
	{
		"0",
		nil,
	},
	{
		"absolutely nit!",
		fmt.Errorf("%q is not a valid answer, please try again.", "absolutely nit!"),
	},
	{
		"huh",
		fmt.Errorf("%q is not a valid answer, please try again.", "huh"),
	},
	{
		"m",
		fmt.Errorf("%q is not a valid answer, please try again.", "m"),
	},
}

func TestYesOrNo(t *testing.T) {
	for _, ynt := range ynTests {
		errString, expString := "nil", "nil"
		err := YesOrNo(ynt.in)

		if err != nil {
			errString = err.Error()
		}

		if ynt.expected != nil {
			expString = ynt.expected.Error()
		}

		if errString != expString {
			t.Errorf("YNT fail have (%s) want (%s)", errString, expString)
		}
	}
}

var numberTests = []struct {
	in       string
	expected error
}{
	{
		"1",
		nil,
	},
	{
		"2.5",
		nil,
	},
	{
		"-3",
		nil,
	},
	{
		"two",
		fmt.Errorf("%q is not a number, please try again.", "two"),
	},
}

func TestIsANumber(t *testing.T) {
	for _, ian := range numberTests {
		errString, expString := "nil", "nil"
		err := IsANumber(ian.in)

		if err != nil {
			errString = err.Error()
		}

		if ian.expected != nil {
			expString = ian.expected.Error()
		}

		if errString != expString {
			t.Errorf("IAN fail have (%s) want (%s)", errString, expString)
		}
	}
}

var graterThanZeroTests = []struct {
	in       string
	expected error
}{
	{
		"1",
		nil,
	},
	{
		"2.5",
		nil,
	},
	{
		"-3",
		fmt.Errorf("%q must be greater than zero, please try again.", "-3"),
	},
	{
		"0",
		fmt.Errorf("%q must be greater than zero, please try again.", "0"),
	},
	{
		"two",
		fmt.Errorf("%q is not a number, please try again.", "two"),
	},
}

func TestGreaterThanZero(t *testing.T) {
	for _, gtz := range graterThanZeroTests {
		errString, expString := "nil", "nil"
		err := GreaterThanZero(gtz.in)

		if err != nil {
			errString = err.Error()
		}

		if gtz.expected != nil {
			expString = gtz.expected.Error()
		}

		if errString != expString {
			t.Errorf("GTZ fail have (%s) want (%s)", errString, expString)
		}
	}
}

var minStringTests = []struct {
	in       string
	min      float64
	expected error
}{
	{
		"blah",
		float64(1),
		nil,
	},
	{
		"blah",
		float64(0),
		nil,
	},
	{
		"blah",
		float64(-1),
		nil,
	},
	{
		"blah",
		float64(5),
		fmt.Errorf("%q must have a min length of %.0f, please try again.", "blah", float64(5)),
	},
	{
		"b",
		float64(1),
		nil,
	},
	{
		"",
		float64(1),
		fmt.Errorf("%q must have a min length of %.0f, please try again.", "", float64(1)),
	},
}

func TestMinString(t *testing.T) {
	for _, mms := range minStringTests {
		errString, expString := "nil", "nil"
		validator := &MinMaxString{Min: mms.min}
		err := validator.CheckMin(mms.in)

		if err != nil {
			errString = err.Error()
		}

		if mms.expected != nil {
			expString = mms.expected.Error()
		}

		if errString != expString {
			t.Errorf("MMS fail have (%s) want (%s)", errString, expString)
		}
	}
}

var maxStringTests = []struct {
	in       string
	max      float64
	expected error
}{
	{
		"blah",
		float64(4),
		nil,
	},
	{
		"blah",
		float64(5),
		nil,
	},
	{
		"blah",
		float64(-1),
		nil,
	},
	{
		"blah",
		float64(3),
		fmt.Errorf("%q must have a max length of %.0f, please try again.", "blah", float64(3)),
	},
	{
		"b",
		float64(1),
		nil,
	},
}

func TestMaxString(t *testing.T) {
	for _, mms := range maxStringTests {
		errString, expString := "nil", "nil"
		validator := &MinMaxString{Max: mms.max}
		err := validator.CheckMax(mms.in)

		if err != nil {
			errString = err.Error()
		}

		if mms.expected != nil {
			expString = mms.expected.Error()
		}

		if errString != expString {
			t.Errorf("MMS fail have (%s) want (%s)", errString, expString)
		}
	}
}

var minNumberTests = []struct {
	in       string
	min      float64
	expected error
}{
	{
		"4",
		float64(1),
		nil,
	},
	{
		"4",
		float64(0),
		nil,
	},
	{
		"4",
		float64(-1),
		nil,
	},
	{
		"4",
		float64(5),
		fmt.Errorf("%q must be greater than %.0f, please try again.", "4", float64(5)),
	},
	{
		"1",
		float64(1),
		nil,
	},
	{
		"0",
		float64(1),
		fmt.Errorf("%q must be greater than %.0f, please try again.", "0", float64(1)),
	},
	{
		"boo",
		float64(10),
		fmt.Errorf("%q is not a number, please try again.", "boo"),
	},
}

func TestMinNumber(t *testing.T) {
	for _, mms := range minNumberTests {
		errString, expString := "nil", "nil"
		validator := &MinMaxNumber{Min: mms.min}
		err := validator.CheckMin(mms.in)

		if err != nil {
			errString = err.Error()
		}

		if mms.expected != nil {
			expString = mms.expected.Error()
		}

		if errString != expString {
			t.Errorf("MMS fail have (%s) want (%s)", errString, expString)
		}
	}
}

var maxNumberTests = []struct {
	in       string
	max      float64
	expected error
}{
	{
		"4",
		float64(4),
		nil,
	},
	{
		"4",
		float64(5),
		nil,
	},
	{
		"4",
		float64(-1),
		nil,
	},
	{
		"4",
		float64(3),
		fmt.Errorf("%q must be less than %.0f, please try again.", "4", float64(3)),
	},
	{
		"1",
		float64(1),
		nil,
	},
	{
		"1",
		float64(10),
		nil,
	},
	{
		"boo",
		float64(10),
		fmt.Errorf("%q is not a number, please try again.", "boo"),
	},
}

func TestMaxNumber(t *testing.T) {
	for _, mms := range maxNumberTests {
		errString, expString := "nil", "nil"
		validator := &MinMaxNumber{Max: mms.max}
		err := validator.CheckMax(mms.in)

		if err != nil {
			errString = err.Error()
		}

		if mms.expected != nil {
			expString = mms.expected.Error()
		}

		if errString != expString {
			t.Errorf("MMS fail have (%s) want (%s)", errString, expString)
		}
	}
}
