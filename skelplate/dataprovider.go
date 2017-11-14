package skelplate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/brainicorn/skelp/prompter"
	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
)

type Flag byte

const (
	UseDefaults Flag = 1 << iota
	SkipMulti
)

const (
	skelpFilename        = "skelp.json"
	hooksDirName         = "hooks"
	ErrSkelpFileNotFound = "skelp.json not found: %s"
	indexSeparator       = ";"
)

type SkelplateDataProvider struct {
	data         map[string]interface{}
	funcMap      map[string]interface{}
	flags        Flag
	tOptions     []string
	beforePrompt func()
}

func NewDataProvider(data map[string]interface{}, opts Flag) *SkelplateDataProvider {
	return &SkelplateDataProvider{
		data:     data,
		funcMap:  skelputil.FunctionMap(),
		tOptions: skelputil.TemplateOptions(),
		flags:    opts,
	}
}

func (sdp *SkelplateDataProvider) DataProviderFunc(templateRoot string) (interface{}, error) {
	var err error
	var data map[string]interface{}
	var descriptorBytes []byte
	var skelplate SkelplateDescriptor

	jsonPath := filepath.Join(templateRoot, skelpFilename)
	if !skelputil.PathExists(jsonPath) {
		err = fmt.Errorf(ErrSkelpFileNotFound, jsonPath)
	}

	if err == nil {
		descriptorBytes, err = ioutil.ReadFile(jsonPath)
	}

	if err == nil {
		skelplate, err = ValidateDescriptor(descriptorBytes)
	}

	if err == nil {
		data, err = sdp.gatherData(skelplate)
	}

	return data, err
}

func (sdp *SkelplateDataProvider) HookProviderFunc(templateRoot string) (provider.Hooks, error) {
	var err error
	var descriptorBytes []byte
	var skelplate SkelplateDescriptor

	hooks := provider.Hooks{}

	jsonPath := filepath.Join(templateRoot, skelpFilename)
	if !skelputil.PathExists(jsonPath) {
		err = fmt.Errorf(ErrSkelpFileNotFound, jsonPath)
	}

	if err == nil {
		descriptorBytes, err = ioutil.ReadFile(jsonPath)
	}

	if err == nil {
		skelplate, err = ValidateDescriptor(descriptorBytes)
	}

	hooksPath := filepath.Join(templateRoot, hooksDirName)

	if err == nil {
		hooks.PreInput, err = convertHooks(hooksPath, skelplate.TemplateHooks.PreInput)
	}

	if err == nil {
		hooks.PreGen, err = convertHooks(hooksPath, skelplate.TemplateHooks.PreGen)
	}

	if err == nil {
		hooks.PostGen, err = convertHooks(hooksPath, skelplate.TemplateHooks.PostGen)
	}

	return hooks, err
}

func convertHooks(basePath string, inHooks []string) ([]string, error) {
	var err error
	outHooks := make([]string, len(inHooks))

	for i, h := range inHooks {
		var absHook string

		if err == nil {
			absHook, err = filepath.Abs(filepath.Join(basePath, h))
			outHooks[i] = absHook
		}
	}

	return outHooks, err
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
	curParents := parents
	for _, v := range variables {
		if err == nil {

			if cv, isCV := v.(*ComplexVar); isCV {
				disabled := false
				disabled, err = sdp.isVariableDisabled(v, fillerData)

				if err == nil {

					if disabled {
						continue
					}

					var askAgain prompter.Prompter
					var answer interface{}
					varname := v.Name()
					curParents = append(curParents, varname)

					if !skelputil.IsBlank(cv.AddPrompt) {
						secondQuestion := cv.AddPrompt

						askAgain = &prompter.KeyedInput{
							Prompt: prompter.Prompt{
								Question: secondQuestion,
								Default:  "y",
							},
							IsConfirm: true,
						}

						objects := make([]map[string]interface{}, 0)
						again := true

						// the resulting data as well as the provided data is an array of maps.
						// we need to make sure we keep track of the array index when prepopulating data
						curIndex := 0
						curParents[len(curParents)-1] = curParents[len(curParents)-1] + indexSeparator + strconv.Itoa(curIndex)

						// we also need to ensure we loop at least as many times as we have provided data
						objPath := strings.Join(curParents, ".")
						minIndex := 0
						if providedObjArray, gotDataObjs := getDotKeyFromMap(objPath, sdp.data); gotDataObjs {
							minIndex = len(providedObjArray.([]map[string]interface{})) - 1
						}

						for again {
							cvData := make(map[string]interface{})
							err = sdp.gatherVariableData(cv.Variables(), cvData, curParents)

							again = (err == nil)
							if err == nil {

								objects = append(objects, cvData)

								if curIndex < minIndex {
									again = true
								} else if sdp.flags&SkipMulti != 0 {
									again = false
								} else {
									if sdp.beforePrompt != nil {
										sdp.beforePrompt()
									}
									again, _ = prompter.AsBool(askAgain.Ask())
								}
								curIndex++
								curParents[len(curParents)-1] = varname + indexSeparator + strconv.Itoa(curIndex)
							}

						}

						answer = objects
					} else {
						cvData := make(map[string]interface{})
						err = sdp.gatherVariableData(cv.Variables(), cvData, curParents)

						answer = cvData
					}

					if err == nil {
						fillerData[varname] = answer
					}
					curParents = []string{}
				}
			} else {
				err = sdp.gatherSingleVariable(v, fillerData, curParents)
			}

		}
	}

	return err
}

func (sdp *SkelplateDataProvider) gatherSingleVariable(v TemplateVariable, fillerData map[string]interface{}, parents []string) error {
	var err error
	var dataval interface{}
	var defval interface{}
	var disabled bool
	prefilled := false

	disabled, err = sdp.isVariableDisabled(v, fillerData)

	if err != nil || disabled {
		return err
	}

	if dv, isDefaultable := v.(Defaultable); isDefaultable {
		defval, err = sdp.renderDefaultValue(dv, fillerData)
	}

	if err != nil {
		return err
	}

	err = sdp.renderChoices(v, fillerData)

	if err != nil {
		return err
	}

	varname := v.Name()

	prefilled, err = sdp.prefillDataVar(varname, defval, fillerData, parents)

	if err != nil {
		return err
	}

	if prefilled {
		return nil
	}

	dataval, err = sdp.promptForVariable(v, varname, defval)

	if err != nil {
		return fmt.Errorf("error asking for input: (%s): %s", varname, err)
	}

	fillerData[varname] = dataval
	return err
}

func (sdp *SkelplateDataProvider) prefillDataVar(varname string, defval interface{}, fillerData map[string]interface{}, parents []string) (bool, error) {
	var dataval interface{}
	var gotdata bool
	var err error
	varpath := strings.Join(append(parents, varname), ".")

	if dataval, gotdata = getDotKeyFromMap(varpath, sdp.data); gotdata {
		fillerVal := dataval

		typeOfDefval := reflect.TypeOf(defval)
		typeOfDataval := reflect.TypeOf(dataval)
		if typeOfDefval.Kind() == typeOfDataval.Kind() {
			if typeOfDataval.Kind() == reflect.String {
				fillerVal, err = sdp.runStringTemplate(dataval.(string), fillerData)

				if err != nil {
					return false, fmt.Errorf("unable to parse data template: %s - %s", dataval, err)
				}
			}

			fillerData[varname] = fillerVal
			return true, nil
		} else {
			return false, fmt.Errorf("invalid type for provided data entry '%s': want (%s) have (%s)", varname, typeOfDataval.Kind(), typeOfDefval.Kind())
		}
	}

	return false, nil
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

func (sdp *SkelplateDataProvider) renderDefaultValue(v Defaultable, fillerData map[string]interface{}) (interface{}, error) {
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
	pathParts := strings.Split(key, ".")
	parent := data
	for i, k := range pathParts {
		key, aryIndex := getKeyAndIndex(k)
		if val, ok := parent[key]; ok {
			if i == len(pathParts)-1 {
				return val, true
			}

			if mapval, ismap := val.(map[string]interface{}); ismap {
				parent = mapval
			} else if arymapval, isarymap := val.([]map[string]interface{}); isarymap && aryIndex > -1 {
				parent = arymapval[aryIndex]
			} else {
				break
			}
		}
	}

	return nil, false
}

func getKeyAndIndex(varname string) (string, int) {
	keyParts := strings.Split(varname, indexSeparator)

	if len(keyParts) == 2 {
		if i, err := strconv.Atoi(keyParts[1]); err == nil {
			return keyParts[0], i
		}
	}
	return varname, -1
}
