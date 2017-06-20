package skelplate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/brainicorn/skelp/prompter"
	"github.com/brainicorn/skelp/skelputil"
)

const (
	promptEnterValue    = "Enter a value for %s:"
	promptMakeSelection = "Make a selection for %s:"
	promptAddAnother    = "Would you like to add another value for %s:"
)

func promptForVariable(tvar TemplateVariable, varname string, dval interface{}, beforePrompt func()) (interface{}, error) {
	var ask prompter.Prompter
	var askAgain prompter.Prompter

	prompt := prompter.Prompt{BeforePrompt: beforePrompt}

	defIsBool := reflect.TypeOf(dval).Kind() == reflect.Bool

	switch tvar.(type) {
	case *SimpleVar:
		ttv := tvar.(*SimpleVar)
		cplxVar := ComplexVar{
			SimpleVar: *ttv,
		}
		configurePrompt(&prompt, cplxVar, varname, promptEnterValue, dval)
		ask = &prompter.KeyedInput{Prompt: prompt, IsConfirm: defIsBool}

	case *ComplexVar:
		ttv := tvar.(*ComplexVar)
		configurePrompt(&prompt, *ttv, varname, promptEnterValue, dval)
		ask = &prompter.KeyedInput{Prompt: prompt, IsConfirm: defIsBool}

	case *MultiValue:
		ttv := tvar.(*MultiValue)
		configurePrompt(&prompt, ttv.ComplexVar, varname, promptEnterValue, dval)
		ask = &prompter.KeyedInput{Prompt: prompt, IsConfirm: false}

		secondQuestion := fmt.Sprintf(promptAddAnother, varname)
		if !skelputil.IsBlank(ttv.AddPrompt) {
			secondQuestion = ttv.AddPrompt
		}

		askAgain = &prompter.KeyedInput{
			Prompt: prompter.Prompt{
				Question: secondQuestion,
				Default:  "y",
			},
			IsConfirm: true,
		}

	case *Selection:
		ttv := tvar.(*Selection)
		configurePrompt(&prompt, ttv.ComplexVar, varname, promptMakeSelection, dval)

		ask = &prompter.SelectedInput{
			Prompt:  prompt,
			Options: ttv.Choices,
			IsMulti: ttv.MultipleChoice,
		}
	}

	return doPrompt(ask, askAgain, beforePrompt, dval)
}

func configurePrompt(prompt *prompter.Prompt, cv ComplexVar, varname, fallbackQuestion string, defval interface{}) {
	prompt.Question = formatQuestion(cv, varname, fallbackQuestion)
	prompt.Validators = []prompter.Validator{}
	configureDefaultAndValidators(prompt, cv, defval)
}

func formatQuestion(cv ComplexVar, varname, fallback string) string {
	question := fmt.Sprintf(fallback, varname)

	if !skelputil.IsBlank(cv.Prompt) {
		question = cv.Prompt
	}

	return question
}

func configureDefaultAndValidators(prompt *prompter.Prompt, cv ComplexVar, defval interface{}) {
	var defstring string

	switch defval.(type) {
	case string:
		defstring = defval.(string)
		if cv.Required {
			prompt.Validators = append(prompt.Validators, prompter.StringNotBlank)
		}

		if cv.Min > 0 || cv.Max > 0 {
			fmt.Println("got min/max")
			mm := &prompter.MinMaxString{
				Min: cv.Min,
				Max: cv.Max,
			}
			prompt.Validators = append(prompt.Validators, mm.CheckMin)
			prompt.Validators = append(prompt.Validators, mm.CheckMax)
			fmt.Println("validators", prompt.Validators)
		}
	case float64:
		defstring = strconv.FormatFloat(defval.(float64), 'f', -1, 64)
		if cv.Required {
			prompt.Validators = append(prompt.Validators, prompter.GreaterThanZero)
		} else {
			prompt.Validators = append(prompt.Validators, prompter.IsANumber)
		}

		if cv.Min > 0 || cv.Max > 0 {
			mm := &prompter.MinMaxNumber{
				Min: cv.Min,
				Max: cv.Max,
			}
			prompt.Validators = append(prompt.Validators, mm.CheckMin)
			prompt.Validators = append(prompt.Validators, mm.CheckMax)
		}
	case bool:
		defstring = strconv.FormatBool(defval.(bool))
	case []interface{}:
		defslice := []string{}
		for _, elem := range defval.([]interface{}) {
			if reflect.TypeOf(elem).Kind() == reflect.String {
				defslice = append(defslice, elem.(string))
			}

			if reflect.TypeOf(elem).Kind() == reflect.Float64 {
				defslice = append(defslice, strconv.FormatFloat(elem.(float64), 'f', -1, 64))
			}
		}
		defstring = strings.Join(defslice, ",")
	}

	prompt.Default = defstring
}

func doPrompt(ask, askAgain prompter.Prompter, beforePrompt func(), defval interface{}) (interface{}, error) {
	var err error
	var answer string
	var finalAnswer interface{}

	if askAgain != nil {
		var ans string
		answers := []string{}
		again := true

		for again {
			ans, err = ask.Ask()
			if err == nil {

				answers = append(answers, ans)

				if beforePrompt != nil {
					beforePrompt()
				}
			}

			again, _ = prompter.AsBool(askAgain.Ask())
		}

		answer = strings.Join(answers, ",")
	} else {
		answer, err = ask.Ask()
	}

	if err == nil {
		finalAnswer, err = convertAnswer(answer, defval)
	}

	return finalAnswer, err
}

func convertAnswer(answer string, defval interface{}) (interface{}, error) {
	var err error
	var typedAnswer interface{}

	switch defval.(type) {
	case string:
		typedAnswer = answer
	case float64:
		typedAnswer, err = strconv.ParseFloat(answer, 64)
	case bool:
		typedAnswer, err = strconv.ParseBool(answer)
	case []interface{}:
		ansSlice := strings.Split(answer, ",")
		var tans interface{}
		typedSlice := []interface{}{}

		switch reflect.TypeOf(defval.([]interface{})[0]).Kind() {
		case reflect.String:
			for _, s := range ansSlice {
				typedSlice = append(typedSlice, s)
			}
		case reflect.Float64:
			for _, s := range ansSlice {
				tans, err = strconv.ParseFloat(s, 64)

				if err == nil {
					typedSlice = append(typedSlice, tans)
				}
			}
		}
		typedAnswer = typedSlice
	}

	return typedAnswer, err
}
