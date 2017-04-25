package schema

const (
	// GithubComBrainicornSkelpTemplateTemplate is a json-schema accessor
	GithubComBrainicornSkelpTemplateTemplate = `{"$schema":"http://json-schema.org/draft-04/schema#","id":"http://github.com/brainicorn/skelp/template/Template","type":"object","definitions":{"github_com-brainicorn-skelp-template-Configurable":{"type":"object","title":"Configurable is an object that can express complex rules for capturing input.","properties":{"required":{"type":"boolean"}}}},"properties":{"author":{"type":"string","title":"TemplateAuthor is the author of the template."},"created":{"type":"string","title":"TemplateCreated is the date the template was created.","format":"date-time"},"modified":{"type":"string","title":"TemplateModified is the date the template was last modified.","format":"date-time"},"repository":{"type":"string","title":"TemplateURL is the url of the template."},"si":{"type":"object","title":"SomeInterface is an interface.","additionalProperties":true},"variables":{"type":"object","title":"TemplateVariables holds the variables and their configuration for processing a template.","additionalProperties":{"type":"object","title":"TemplateDataTypes are the types allowed for variables in templates.","anyOf":[{"type":"string"},{"type":"number"},{"type":"integer"},{"type":"boolean"},{"type":"array"},{"$ref":"#/definitions/github_com-brainicorn-skelp-template-Configurable"}]}}}}`

	// GithubComBrainicornSkelpTemplateConfigurable is a json-schema accessor
	GithubComBrainicornSkelpTemplateConfigurable = `{"$schema":"http://json-schema.org/draft-04/schema#","id":"http://github.com/brainicorn/skelp/template/Configurable","type":"object","title":"Configurable is an object that can express complex rules for capturing input.","properties":{"required":{"type":"boolean"}}}`

)
