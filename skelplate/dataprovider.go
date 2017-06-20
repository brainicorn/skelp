package skelplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/brainicorn/skelp/skelputil"
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

	jsonPath := filepath.Join(templateRoot, skelpFilename)
	if !skelputil.PathExists(jsonPath) {
		err = fmt.Errorf(ErrSkelpFileNotFound, jsonPath)
	}

	if err == nil {
		descriptorBytes, err = ioutil.ReadFile(jsonPath)
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

	for _, v := range descriptor.TemplateVariables {
		var dataval interface{}
		var gotdata bool
		var defval interface{}
		varname, err := sdp.runStringTemplate(v.Name(), fillerData)

		if err != nil {
			return nil, fmt.Errorf("unable to parse variable name template: %s - %s", v.Name(), err)
		}

		valOfDefault := reflect.ValueOf(v.Default())

		if valOfDefault.Kind() == reflect.String {
			defval, err = sdp.runStringTemplate(v.Default().(string), fillerData)

			if err != nil {
				return nil, fmt.Errorf("unable to parse variable default template: %s - %s", v.Default(), err)
			}
		} else if isStringSlice(valOfDefault) {
			defOpts := v.Default().([]interface{})
			defVals := []interface{}{}
			for _, ds := range defOpts {
				dv, dverr := sdp.runStringTemplate(ds.(string), fillerData)
				if dverr != nil {
					return nil, fmt.Errorf("unable to parse variable default template: %s - %s", ds, dverr)
				}
				defVals = append(defVals, dv)
			}
			defval = defVals

		} else {
			defval = v.Default()
		}

		if dataval, gotdata = sdp.data[varname]; gotdata {
			fillerVal := dataval
			typeOfDefval := reflect.TypeOf(defval)
			typeOfDataval := reflect.TypeOf(dataval)
			if typeOfDefval.Kind() == typeOfDataval.Kind() {
				if typeOfDataval.Kind() == reflect.String {
					fillerVal, err = sdp.runStringTemplate(dataval.(string), fillerData)

					if err != nil {
						return nil, fmt.Errorf("unable to parse data template: %s - %s", dataval, err)
					}
				}

				fillerData[varname] = fillerVal
				continue
			} else {
				return nil, fmt.Errorf("invalid type for provided data entry '%s': want (%s) have (%s)", varname, typeOfDataval.Kind(), typeOfDefval.Kind())
			}
		}

		dataval, err = promptForVariable(v, varname, defval, sdp.beforePrompt)

		if err != nil {
			return nil, fmt.Errorf("error asking for input: (%s): %s", varname, err)
		}

		fillerData[varname] = dataval

	}

	return fillerData, err
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
