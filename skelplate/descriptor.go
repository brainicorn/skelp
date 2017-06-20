package skelplate

import (
	"encoding/json"
	"time"
)

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

// TemplateVariable is the base interface for a variable
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/skelp/skelplate/SimpleVar"
//	,"github.com/brainicorn/skelp/skelplate/ComplexVar"
//	,"github.com/brainicorn/skelp/skelplate/Selection"
//	,"github.com/brainicorn/skelp/skelplate/MultiValue"]
// )
type TemplateVariable interface {
	Name() string
	Default() interface{}
}

// SimpleVar is an object that can express a name value pair
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
	// @jsonSchema(required=true, anyOf=["string","number","integer","boolean","array"])
	DefaultVal interface{} `json:"default"`
}

func (sv *SimpleVar) Name() string {
	return sv.Varname
}

func (sv *SimpleVar) Default() interface{} {
	return sv.DefaultVal
}

// TODO add a Keyed struct and remove min/max from complex.
// Complex becomes just a base type. Keyed should also have a IsPassword field.

// ComplexVar is an object that can express complex rules for capturing input.
//
// @jsonSchema(additionalProperties=false)
type ComplexVar struct {
	SimpleVar

	// Required whether or not a non-empty value is required.
	Required bool `json:"required,omitempty"`

	// Prompt the string to display when asking for a value.
	Prompt string `json:"prompt,omitempty"`

	// Min the minimum value (for numbers) or length (for strings).
	Min float64 `json:"min,omitempty"`

	// Max the maximum value (for numbers) or length (for strings)
	Max float64 `json:"max,omitempty"`
}

// Selection represents a configurable "select box".
// The user can choose multiple values or be restricted to choosing a single value.
//
// @jsonSchema(additionalProperties=false)
type Selection struct {
	ComplexVar

	// MultipleChoice designates whether multiple valuse may be chosen when the choices field is present.
	//
	// @jsonSchema(required=true)
	MultipleChoice bool `json:"mutlichoice"`

	// Choices are the options to display in a select box.
	// @jsonSchema(required=true)
	Choices []string `json:"choices,omitempty"`
}

// MultiValue allows the user to enter multiple values.
// This is for gathering things like "tags"
//
// @jsonSchema(additionalProperties=false)
type MultiValue struct {
	ComplexVar

	// IsMultiVal designates the variable as a mutli-value prompt.
	//
	// @jsonSchema(required=true)
	IsMultiVal bool `json:"mutlival"`

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
						if _, selok := vvmap["choices"]; selok {
							var sel Selection
							var seljs []byte
							seljs, err = json.Marshal(vv)
							if err == nil {
								err = json.Unmarshal(seljs, &sel)
								if err == nil {
									varSlice = append(varSlice, &sel)
								}
							}
						} else if _, mvok := vvmap["mutlival"]; mvok {
							var mv MultiValue
							var mvjs []byte
							mvjs, err = json.Marshal(vv)
							if err == nil {
								err = json.Unmarshal(mvjs, &mv)
								if err == nil {
									varSlice = append(varSlice, &mv)
								}
							}
						} else if len(vvmap) > 2 {
							var cplx ComplexVar
							var cplxjs []byte
							cplxjs, err = json.Marshal(vv)
							if err == nil {
								err = json.Unmarshal(cplxjs, &cplx)
								if err == nil {
									varSlice = append(varSlice, &cplx)
								}
							}
						} else {
							var smp SimpleVar
							var smpjs []byte
							smpjs, err = json.Marshal(vv)
							if err == nil {
								err = json.Unmarshal(smpjs, &smp)
								if err == nil {
									varSlice = append(varSlice, &smp)
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
