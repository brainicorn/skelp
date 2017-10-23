package skelplate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/brainicorn/skelp/prompter"
	"github.com/brainicorn/skelp/skelputil"
	"github.com/xeipuuv/gojsonschema"
)

const (
	skelpFilename        = "skelp.json"
	ErrSkelpFileNotFound = "skelp.json not found: %s"
)

type SkelplateDataProvider struct {
	data         map[string]interface{}
	funcMap      map[string]interface{}
	tOptions     []string
	beforePrompt func()
}

func NewDataProvider(data map[string]interface{}) *SkelplateDataProvider {
	return &SkelplateDataProvider{
		data:     data,
		funcMap:  skelputil.FunctionMap(),
		tOptions: skelputil.TemplateOptions(),
	}
}

func (sdp *SkelplateDataProvider) DataProviderFunc(templateRoot string) (interface{}, error) {
	var err error
	var data map[string]interface{}
	var descriptorBytes []byte
	var skelplate SkelplateDescriptor
	var schemaValidationResult *gojsonschema.Result
	jsonPath := filepath.Join(templateRoot, skelpFilename)
	if !skelputil.PathExists(jsonPath) {
		err = fmt.Errorf(ErrSkelpFileNotFound, jsonPath)
	}

	if err == nil {
		descriptorBytes, err = ioutil.ReadFile(jsonPath)
	}

	if err == nil {
		schemaLoader := gojsonschema.NewStringLoader(GithubComBrainicornSkelpSkelplateSkelplateDescriptor)
		docLoader := gojsonschema.NewBytesLoader(descriptorBytes)

		schemaValidationResult, err = gojsonschema.Validate(schemaLoader, docLoader)

		if err == nil && len(schemaValidationResult.Errors()) > 0 {
			var errBuf bytes.Buffer
			errBuf.WriteString("Error validating skelp descriptor:\n")
			for _, re := range schemaValidationResult.Errors() {
				errBuf.WriteString(fmt.Sprintf("  - %s\n", re))
			}

			err = errors.New(errBuf.String())
		}
	}
	if err == nil {
		err = json.Unmarshal(descriptorBytes, &skelplate)
	}
	if err == nil {
		data, err = sdp.gatherData(skelplate)
	}

	return data, err
}

func (sdp *SkelplateDataProvider) gatherData(descriptor SkelplateDescriptor) (map[string]interface{}, error) {
	var err error
	fillerData := make(map[string]interface{})

	fillerData["TemplateAuthor"] = descriptor.TemplateAuthor
	fillerData["TemplateRepo"] = descriptor.TemplateRepo
	fillerData["TemplateCreated"] = descriptor.TemplateCreated
	fillerData["TemplateModified"] = descriptor.TemplateModified
	fillerData["TemplateDesc"] = descriptor.TemplateDesc

	err = sdp.gatherVariableData(descriptor.Variables(), fillerData, []string{})

	return fillerData, err
}

func (sdp *SkelplateDataProvider) gatherVariableData(variables []TemplateVariable, fillerData map[string]interface{}, parents []string) error {
	var err error
	for _, v := range variables {
		if err == nil {
			if cv, isCV := v.(*ComplexVar); isCV {
				disabled := false
				disabled, err = sdp.isVariableDisabled(v, fillerData)

				if err != nil || disabled {
					return err
				}

				var askAgain prompter.Prompter
				var answer interface{}
				varname := v.Name()
				parents = append(parents, varname)

				if cv.IsMultiVal {
					secondQuestion := fmt.Sprintf(promptAddAnother, varname)
					if !skelputil.IsBlank(cv.AddPrompt) {
						secondQuestion = cv.AddPrompt
					}

					askAgain = &prompter.KeyedInput{
						Prompt: prompter.Prompt{
							Question: secondQuestion,
							Default:  "y",
						},
						IsConfirm: true,
					}
				}

				if askAgain != nil {
					objects := make([]map[string]interface{}, 0)
					again := true

					for again {
						cvData := make(map[string]interface{})
						err = sdp.gatherVariableData(cv.Variables(), cvData, parents)
						if err == nil {

							objects = append(objects, cvData)

							if sdp.beforePrompt != nil {
								sdp.beforePrompt()
							}
						}

						again, _ = prompter.AsBool(askAgain.Ask())
					}

					answer = objects
				} else {
					cvData := make(map[string]interface{})
					err = sdp.gatherVariableData(cv.Variables(), cvData, parents)

					answer = cvData
				}

				if err == nil {
					fillerData[varname] = answer
				}
			} else {
				err = sdp.gatherSingleVariable(v, fillerData, parents)
			}
		}
	}

	return err
}

func (sdp *SkelplateDataProvider) gatherSingleVariable(v TemplateVariable, fillerData map[string]interface{}, parents []string) error {
	var err error
	var dataval interface{}
	var gotdata bool
	var defval interface{}
	var disabled bool
	disabled, err = sdp.isVariableDisabled(v, fillerData)

	if err != nil || disabled {
		return err
	}

	defval, err = sdp.renderDefaultValue(v, fillerData)

	if err != nil {
		return err
	}

	err = sdp.renderChoices(v, fillerData)

	if err != nil {
		return err
	}

	varname := v.Name()

	varpath := strings.Join(append(parents, varname), ".")

	// todo write a helper to get/set from nested maps using dot notation
	if dataval, gotdata = getDotKeyFromMap(varpath, sdp.data); gotdata {
		fillerVal := dataval
		typeOfDefval := reflect.TypeOf(defval)
		typeOfDataval := reflect.TypeOf(dataval)
		if typeOfDefval.Kind() == typeOfDataval.Kind() {
			if typeOfDataval.Kind() == reflect.String {
				fillerVal, err = sdp.runStringTemplate(dataval.(string), fillerData)

				if err != nil {
					return fmt.Errorf("unable to parse data template: %s - %s", dataval, err)
				}
			}

			fillerData[varname] = fillerVal
			return nil
		} else {
			return fmt.Errorf("invalid type for provided data entry '%s': want (%s) have (%s)", varname, typeOfDataval.Kind(), typeOfDefval.Kind())
		}
	}

	dataval, err = sdp.promptForVariable(v, varname, defval)

	if err != nil {
		return fmt.Errorf("error asking for input: (%s): %s", varname, err)
	}

	fillerData[varname] = dataval
	return err
}

func (sdp *SkelplateDataProvider) isVariableDisabled(v TemplateVariable, fillerData map[string]interface{}) (bool, error) {
	var err error
	var disabledString string
	var disabled bool

	disabledString, err = sdp.runStringTemplate(v.DisabledTemplate(), fillerData)
	if err != nil {
		return false, fmt.Errorf("unable to parse variable disabled template: %s - %s", v.DisabledTemplate(), err)
	}

	disabled, err = strconv.ParseBool(disabledString)

	if err != nil {
		return false, fmt.Errorf("unable to convert disabled value to a boolean: %s - %s", disabledString, err)
	}

	return disabled, nil
}

func (sdp *SkelplateDataProvider) renderDefaultValue(v TemplateVariable, fillerData map[string]interface{}) (interface{}, error) {
	var err error
	var defval interface{}

	valOfDefault := reflect.ValueOf(v.Default())

	if valOfDefault.Kind() == reflect.String {
		defval, err = sdp.runStringTemplate(v.Default().(string), fillerData)

		if err != nil {
			err = fmt.Errorf("unable to parse variable default template: %s - %s", v.Default(), err)
		}
	} else if isStringSlice(valOfDefault) {
		defOpts := v.Default().([]interface{})
		defVals := []interface{}{}
		for _, ds := range defOpts {
			dv, dverr := sdp.runStringTemplate(ds.(string), fillerData)
			if dverr != nil {
				err = fmt.Errorf("unable to parse variable default template: %s - %s", ds, dverr)
				break
			}
			defVals = append(defVals, dv)
		}
		defval = defVals
	} else {
		defval = v.Default()
	}

	return defval, err
}

func (sdp *SkelplateDataProvider) renderChoices(v TemplateVariable, fillerData map[string]interface{}) error {
	var err error
	var choices = []interface{}{}

	if sel, isSel := v.(*Selection); isSel {
		if reflect.TypeOf(sel.Choices[0]).Kind() == reflect.String {
			for _, choiceTemplate := range sel.Choices {
				var choice string
				choice, err = sdp.runStringTemplate(choiceTemplate.(string), fillerData)

				if err != nil {
					return fmt.Errorf("unable to parse choice template: %s - %s", choiceTemplate, err)
				}

				if len(strings.TrimSpace(choice)) > 0 {
					choices = append(choices, choice)
				}
			}
			sel.Choices = choices
		}
	}

	return nil
}

func (sdp *SkelplateDataProvider) runStringTemplate(input string, tmplData interface{}) (string, error) {
	var err error
	var target string
	var inputTmpl *template.Template
	var b bytes.Buffer

	if !strings.Contains(input, "{{") {
		return input, nil
	}

	if err == nil {
		inputTmpl, err = template.New("nameOrDefault template").Option(sdp.tOptions...).Funcs(sdp.funcMap).Parse(input)
	}

	if err == nil {
		err = inputTmpl.Execute(&b, tmplData)
	}

	if err == nil {
		target = b.String()
	}

	return target, err
}

func isStringSlice(valOf reflect.Value) bool {
	if valOf.Kind() == reflect.Slice {
		if valOf.Len() > 0 {
			_, isstr := valOf.Index(0).Interface().(string)
			return isstr
		}
	}
	return false
}

func getDotKeyFromMap(key string, data map[string]interface{}) (interface{}, bool) {
	parts := strings.Split(key, ".")
	parent := data
	for i, k := range parts {
		if val, ok := parent[k]; ok {
			if i == len(parts)-1 {
				return val, true
			}

			if mapval, ismap := val.(map[string]interface{}); ismap {
				parent = mapval
			} else {
				break
			}
		}
	}

	return nil, false
}
