package prompter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	truePattern  = regexp.MustCompile("^(?i:y(?:es)?)|(?i:t(?:rue)?)|1$")
	falsePattern = regexp.MustCompile("^(?i:n(?:o)?)|(?i:f(?:alse)?)|0$")
)

type Validator func(val string) error

func StringNotBlank(val string) error {
	if len(strings.TrimSpace(val)) < 1 {
		return fmt.Errorf("a value is required")
	}
	return nil
}

func YesOrNo(val string) error {
	if !truePattern.Match([]byte(val)) && !falsePattern.Match([]byte(val)) {
		return fmt.Errorf("%q is not a valid answer, please try again.", val)
	}
	return nil
}

func IsANumber(val string) error {
	_, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return fmt.Errorf("%q is not a number, please try again.", val)
	}

	return nil
}

func GreaterThanZero(val string) error {
	f, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return fmt.Errorf("%q is not a number, please try again.", val)
	}

	if f < 1 {
		return fmt.Errorf("%q must be greater than zero, please try again.", val)
	}

	return nil
}

type MinMaxString struct {
	Min float64
	Max float64
}

func (mm *MinMaxString) CheckMin(val string) error {
	fmt.Println("min string")
	if mm.Min < 1 {
		return nil
	}

	if float64(len(strings.TrimSpace(val))) < mm.Min {
		return fmt.Errorf("%q must have a min length of %d, please try again.", val, mm.Min)
	}

	return nil
}

func (mm *MinMaxString) CheckMax(val string) error {
	if mm.Max < 1 {
		return nil
	}

	if float64(len(strings.TrimSpace(val))) > mm.Max {
		return fmt.Errorf("%q must have a max length of %d, please try again.", val, mm.Max)
	}

	return nil
}

type MinMaxNumber struct {
	Min float64
	Max float64
}

func (mm *MinMaxNumber) CheckMin(val string) error {
	if mm.Min < 1 {
		return nil
	}

	f, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return fmt.Errorf("%q is not a number, please try again.", val)
	}

	if f < mm.Min {
		return fmt.Errorf("%q must be greater than %d, please try again.", val, mm.Min)
	}

	return nil
}

func (mm *MinMaxNumber) CheckMax(val string) error {
	if mm.Max < 1 {
		return nil
	}

	f, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return fmt.Errorf("%q is not a number, please try again.", val)
	}

	if f > mm.Max {
		return fmt.Errorf("%q must be less than %d, please try again.", val, mm.Min)
	}

	return nil
}
