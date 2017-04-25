package skelplate

import (
	"encoding/json"
	"time"
)

//go:generate jsonschemagen -v -x -c -r -o ./schema github.com/brainicorn/skelp/template Template
type Template struct {
	// TemplateAuthor is the author of the template.
	TemplateAuthor string `json:"author"`

	// TemplateURL is the url of the template.
	TemplateRepo string `json:"repository"`

	// TemplateCreated is the date the template was created.
	TemplateCreated time.Time `json:"created"`

	// TemplateModified is the date the template was last modified.
	TemplateModified time.Time `json:"modified"`

	// TemplateVariables holds the variables and their configuration for processing a template.
	TemplateVariables TemplateData `json:"variables"`
}

// TemplateData is an interface for defining the various types of template inputs.
//
// @jsonSchema(
// 	additionalProperties=["github.com/brainicorn/skelp/template/TemplateDataTypes"]
// )
type TemplateData interface {
}

// UnmarshalJSON converts this bool or schema object from a JSON structure
func (td *Template) UnmarshalJSON(data []byte) error {
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
				td.TemplateCreated = v.(time.Time)
			case "modified":
				td.TemplateModified = v.(time.Time)
			case "variables":
				varMap := make(map[string]interface{})
				vars := v.(map[string]interface{})
				for vk, vv := range vars {
					//it's an object
					if _, ok := vv.(map[string]interface{}); ok {
						var cfg Configurable
						var cfgjs []byte
						cfgjs, err = json.Marshal(vv)
						if err == nil {
							err = json.Unmarshal(cfgjs, &cfg)
							if err == nil {
								varMap[vk] = cfg
							}
						}
					} else {
						varMap[vk] = vv
					}
				}
				td.TemplateVariables = varMap
			}
		}
	}

	return err
}

// TemplateDataTypes are the types allowed for variables in templates.
//
// @jsonSchema(
// 	anyOf=["string","number","integer","boolean","array","github.com/brainicorn/skelp/template/Configurable"]
// )
type TemplateDataTypes interface{}

// Configurable is an object that can express complex rules for capturing input.
type Configurable struct {
	Required bool        `json:"required"`
	Default  interface{} `json:"default"`
	//	AllowMutliple bool        `json:"allowMultiple"`
	//	Prompt        string      `json:"prompt"`
	//	Min           int         `json:"min"`
	//	Max           int         `json:"max"`
	//	Pattern       string      `json:"pattern"`
}
