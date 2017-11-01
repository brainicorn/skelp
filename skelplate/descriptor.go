package skelplate

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	typeSimple     = "simple"
	typeMultiVal   = "multival"
	typeCustomized = "customized"
	typeComplex    = "complex"
	typeSelect     = "select"
)

type VarCollection interface {
	Variables() []TemplateVariable
}

type SkelplateDescriptor struct {
	// TemplateAuthor is the author of the template.
	TemplateAuthor string `json:"author"`

	// TemplateRepo is the url of the template.
	TemplateRepo string `json:"repository"`

	// TemplateDesc is the description of the template.
	TemplateDesc string `json:"description"`

	// TemplateCreated is the date the template was created.
	TemplateCreated time.Time `json:"created"`

	// TemplateModified is the date the template was last modified.
	TemplateModified time.Time `json:"modified"`

	// TemplateVariables holds the variables and their configuration for processing a template.
	TemplateVariables []TemplateVariable `json:"variables"`
}

func (sd *SkelplateDescriptor) Variables() []TemplateVariable {
	return sd.TemplateVariables
}

// TemplateVariable is the base interface for a variable.
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/skelp/skelplate/SimpleVar"
//  ,"github.com/brainicorn/skelp/skelplate/ComplexVar"
//	,"github.com/brainicorn/skelp/skelplate/CustomizedVar"
//	,"github.com/brainicorn/skelp/skelplate/Selection"
//	,"github.com/brainicorn/skelp/skelplate/MultiValue"]
// )
type TemplateVariable interface {
	Name() string
	Default() interface{}
	DisabledTemplate() string
}

// SimpleVar is an object that can express a name value pair.
//
// @jsonSchema(additionalProperties=false)
type SimpleVar struct {

	// Name is the name of the variable.
	// The name can be a golang template and can use values gathered from previous
	// variables in the variables array.
	//
	// @jsonSchema(required=true)
	Varname string `json:"name,omitempty"`

	// Default the default value (can be blank).
	//
	// @jsonSchema(required=true, type=["string","number","integer","boolean","array"])
	DefaultVal interface{} `json:"default"`

	// Disabled will disable this prompt if set to true.
	// NOTE: this value should be a boolean value wrapped in quotes. We've set this aa a string
	// so that golang templating can be used and still pass schema validation.
	// This can be used to make some prompts dependent on other prompts.
	// For example you could set: "disabled":{{ not .SomeOtherBooleanVar }}
	//
	Disabled string `json:"disabled"`
}

func (sv *SimpleVar) Name() string {
	return sv.Varname
}

func (sv *SimpleVar) Default() interface{} {
	return sv.DefaultVal
}

func (sv *SimpleVar) DisabledTemplate() string {
	return sv.Disabled
}

// ComplexVar is an object container for other variables.
//
// @jsonSchema(additionalProperties=false)
type ComplexVar struct {

	// Name is the name of the variable.
	// The name can be a golang template and can use values gathered from previous
	// variables in the variables array.
	//
	// @jsonSchema(required=true)
	Varname string `json:"name,omitempty"`

	// Disabled will disable this prompt if set to true.
	// NOTE: this value should be a boolean value wrapped in quotes. We've set this aa a string
	// so that golang templating can be used and still pass schema validation.
	// This can be used to make some prompts dependent on other prompts.
	// For example you could set: "disabled":{{ not .SomeOtherBooleanVar }}
	//
	Disabled string `json:"disabled"`

	// TemplateVariables holds the variables that make up the fields of the object.
	//
	// @jsonSchema(required=true)
	TemplateVariables []TemplateVariable `json:"variables"`

	// AddPrompt is the string to display when asking if another value should be entered.
	AddPrompt string `json:"addPrompt,omitempty"`
}

func (cv *ComplexVar) Name() string {
	return cv.Varname
}

func (cv *ComplexVar) Default() interface{} {
	return struct{}{}
}

func (cv *ComplexVar) DisabledTemplate() string {
	return cv.Disabled
}

func (cv *ComplexVar) Variables() []TemplateVariable {
	return cv.TemplateVariables
}

// CustomizedVar customizes input.
//
// @jsonSchema(additionalProperties=false)
type CustomizedVar struct {
	SimpleVar

	// Required whether or not a non-empty value is required.
	Required bool `json:"required,omitempty"`

	// Prompt the string to display when asking for a value.
	Prompt string `json:"prompt,omitempty"`

	// Min the minimum value (for numbers) or length (for strings).
	Min float64 `json:"min,omitempty"`

	// Max the maximum value (for numbers) or length (for strings)
	Max float64 `json:"max,omitempty"`

	// Password is a flag to turn on input masking for hiding passwords
	Password bool `json:"password"`
}

// Selection represents a configurable "select box".
// The user can choose multiple values or be restricted to choosing a single value.
//
// @jsonSchema(additionalProperties=false)
type Selection struct {
	CustomizedVar

	// Choices are the options to display in a select box.
	// @jsonSchema(required=true, type=["string","number","integer"])
	Choices []interface{} `json:"choices,omitempty"`
}

// MultiValue allows the user to enter multiple values.
// This is for gathering things like "tags"
//
// @jsonSchema(additionalProperties=false)
type MultiValue struct {
	CustomizedVar

	// AddPrompt is the string to display when asking if another value should be entered.
	AddPrompt string `json:"addPrompt,omitempty"`
}

// UnmarshalJSON cretaes a template object from a JSON structure
func (td *SkelplateDescriptor) UnmarshalJSON(data []byte) error {
	var err error
	var stuff map[string]interface{}

	err = json.Unmarshal(data, &stuff)

	if err == nil {
		for k, v := range stuff {
			switch k {
			case "author":
				td.TemplateAuthor = v.(string)
			case "repository":
				td.TemplateRepo = v.(string)
			case "created":
				td.TemplateCreated, _ = time.Parse(time.RFC3339Nano, v.(string))
			case "modified":
				td.TemplateModified, _ = time.Parse(time.RFC3339Nano, v.(string))
			case "variables":
				varSlice := []TemplateVariable{}
				vars := v.([]interface{})
				for _, vv := range vars {
					//it's an object
					if vvmap, ok := vv.(map[string]interface{}); ok {
						var jsbytes []byte
						jsbytes, err = json.Marshal(vv)
						if err == nil {
							switch typeOfVar(vvmap) {
							case typeSelect:
								var typedVar Selection
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeComplex:
								var typedVar ComplexVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeMultiVal:
								var typedVar MultiValue
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeCustomized:
								var typedVar CustomizedVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeSimple:
								var typedVar SimpleVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							}
						}

					}
				}
				td.TemplateVariables = varSlice
			}
		}
	}

	return err
}

// UnmarshalJSON cretaes a template object from a JSON structure
func (cv *ComplexVar) UnmarshalJSON(data []byte) error {
	var err error
	var stuff map[string]interface{}

	err = json.Unmarshal(data, &stuff)

	if err == nil {
		for k, v := range stuff {
			switch k {
			case "name":
				cv.Varname = v.(string)
			case "disabled":
				cv.Disabled = v.(string)
			case "addPrompt":
				cv.AddPrompt = v.(string)
			case "variables":
				varSlice := []TemplateVariable{}
				vars := v.([]interface{})
				for _, vv := range vars {
					//it's an object
					if vvmap, ok := vv.(map[string]interface{}); ok {
						var jsbytes []byte
						jsbytes, err = json.Marshal(vv)
						if err == nil {
							switch typeOfVar(vvmap) {
							case typeSelect:
								var typedVar Selection
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeComplex:
								var typedVar ComplexVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeMultiVal:
								var typedVar MultiValue
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeCustomized:
								var typedVar CustomizedVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							case typeSimple:
								var typedVar SimpleVar
								err = json.Unmarshal(jsbytes, &typedVar)
								if err == nil {
									if len(strings.TrimSpace(typedVar.Disabled)) < 1 {
										typedVar.Disabled = "false"
									}
									varSlice = append(varSlice, &typedVar)
								}
							}
						}

					}
				}
				cv.TemplateVariables = varSlice
			}
		}
	}

	return err
}

func typeOfVar(varmap map[string]interface{}) string {
	if _, ok := varmap["choices"]; ok {
		return typeSelect
	}

	if _, ok := varmap["variables"]; ok {
		return typeComplex
	}

	if _, ok := varmap["addPrompt"]; ok {
		return typeMultiVal
	}

	rkeys := []string{"min", "max", "password", "prompt", "required"}
	for _, k := range rkeys {
		if _, ok := varmap[k]; ok {
			return typeCustomized
		}
	}

	return typeSimple

}
