package skelplate

const (
	// GithubComBrainicornSkelpSkelplateSkelplateDescriptor is a json-schema accessor
	GithubComBrainicornSkelpSkelplateSkelplateDescriptor = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","definitions":{"github_com-brainicorn-skelp-skelplate-ComplexVar":{"type":"object","title":"ComplexVar is an object container for other variables.","properties":{"addPrompt":{"type":"string","title":"AddPrompt is the string to display when asking if another value should be entered."},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."},"variables":{"type":"array","title":"TemplateVariables holds the variables that make up the fields of the object.","items":{"type":"object","title":"TemplateVariable is the base interface for a variable.","anyOf":[{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-SimpleVar"},{"$ref":"#"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-CustomizedVar"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-Selection"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-MultiValue"}]}}},"required":["name","variables"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-CustomizedVar":{"type":"object","title":"CustomizedVar customizes input.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-Hooks":{"type":"object","title":"Hooks is the object that holds arrays of the various hook scripts.","description":"Each lifecycle can have an array of strings that represent the shell scripts to run.\nEach string should be the basename of the script file followed by any arguments.\nThe string will be processed as a Go Template so the args can use built-in functions and any data\nthat's available from gathering input. The script is assumed to live in the template repo's hooks driectory.\n","properties":{"postGen":{"type":"array","items":{"type":"string"}},"preGen":{"type":"array","items":{"type":"string"}},"preInput":{"type":"array","items":{"type":"string"}}}},"github_com-brainicorn-skelp-skelplate-MultiValue":{"type":"object","title":"MultiValue allows the user to enter multiple values.","description":"This is for gathering things like \"tags\"","properties":{"addPrompt":{"type":"string","title":"AddPrompt is the string to display when asking if another value should be entered."},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-Selection":{"type":"object","title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","properties":{"choices":{"type":"array","title":"Choices are the options to display in a select box.","items":{"type":["string","number","integer"],"title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","additionalProperties":false}},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default","choices"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-SimpleVar":{"type":"object","title":"SimpleVar is an object that can express a name value pair.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."}},"required":["name","default"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-TemplateExclude":{"type":"object","title":"TemplateExclude holds a condition, which if true excludes the list of template paths from processing","properties":{"exclude":{"type":"string","title":"Exclude is a go template that should evalutate to true or false.","description":"If it evaluates to true, the list of template paths will be excluded from processing\n"},"paths":{"type":"array","title":"FilesOrDirs holds the paths that should be excluded when Excludes is true.","items":{"type":"string"}}}}},"properties":{"author":{"type":"string","title":"TemplateAuthor is the author of the template."},"created":{"type":"string","title":"TemplateCreated is the date the template was created.","format":"date-time"},"description":{"type":"string","title":"TemplateDesc is the description of the template."},"excludes":{"type":"array","title":"TemplateExcludes allows for conditionally excluding certain template files from processing","items":{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-TemplateExclude"}},"hooks":{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-Hooks","title":"TemplateHooks holds the scripts that can run during the generation process","description":"Each lifecycle can have an array of strings that represent the shell scripts to run.\nEach string should be the basename of the script file followed by any arguments.\nThe string will be processed as a Go Template so the args can use built-in functions and any data\nthat's available from gathering input. The script is assumed to live in the template repo's hooks driectory.\n"},"modified":{"type":"string","title":"TemplateModified is the date the template was last modified.","format":"date-time"},"repository":{"type":"string","title":"TemplateRepo is the url of the template."},"variables":{"type":"array","title":"TemplateVariables holds the variables and their configuration for processing a template.","items":{"type":"object","title":"TemplateVariable is the base interface for a variable.","anyOf":[{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-SimpleVar"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-ComplexVar"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-CustomizedVar"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-Selection"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-MultiValue"}]}}}}`

	// GithubComBrainicornSkelpSkelplateSimpleVar is a json-schema accessor
	GithubComBrainicornSkelpSkelplateSimpleVar = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"SimpleVar is an object that can express a name value pair.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."}},"required":["name","default"],"additionalProperties":false}`

	// GithubComBrainicornSkelpSkelplateCustomizedVar is a json-schema accessor
	GithubComBrainicornSkelpSkelplateCustomizedVar = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"CustomizedVar customizes input.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false}`

	// GithubComBrainicornSkelpSkelplateSelection is a json-schema accessor
	GithubComBrainicornSkelpSkelplateSelection = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","properties":{"choices":{"type":"array","title":"Choices are the options to display in a select box.","items":{"type":["string","number","integer"],"title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","additionalProperties":false}},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default","choices"],"additionalProperties":false}`

	// GithubComBrainicornSkelpSkelplateMultiValue is a json-schema accessor
	GithubComBrainicornSkelpSkelplateMultiValue = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"MultiValue allows the user to enter multiple values.","description":"This is for gathering things like \"tags\"","properties":{"addPrompt":{"type":"string","title":"AddPrompt is the string to display when asking if another value should be entered."},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false}`

	// GithubComBrainicornSkelpSkelplateComplexVar is a json-schema accessor
	GithubComBrainicornSkelpSkelplateComplexVar = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"ComplexVar is an object container for other variables.","definitions":{"github_com-brainicorn-skelp-skelplate-CustomizedVar":{"type":"object","title":"CustomizedVar customizes input.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-MultiValue":{"type":"object","title":"MultiValue allows the user to enter multiple values.","description":"This is for gathering things like \"tags\"","properties":{"addPrompt":{"type":"string","title":"AddPrompt is the string to display when asking if another value should be entered."},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-Selection":{"type":"object","title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","properties":{"choices":{"type":"array","title":"Choices are the options to display in a select box.","items":{"type":["string","number","integer"],"title":"Selection represents a configurable \"select box\".","description":"The user can choose multiple values or be restricted to choosing a single value.","additionalProperties":false}},"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"max":{"type":"number","title":"Max the maximum value (for numbers) or length (for strings)"},"min":{"type":"number","title":"Min the minimum value (for numbers) or length (for strings)."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"password":{"type":"boolean","title":"Password is a flag to turn on input masking for hiding passwords"},"prompt":{"type":"string","title":"Prompt the string to display when asking for a value."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."}},"required":["name","default","choices"],"additionalProperties":false},"github_com-brainicorn-skelp-skelplate-SimpleVar":{"type":"object","title":"SimpleVar is an object that can express a name value pair.","properties":{"default":{"type":["string","number","integer","boolean","array"],"title":"Default the default value (can be blank).","additionalProperties":false},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."}},"required":["name","default"],"additionalProperties":false}},"properties":{"addPrompt":{"type":"string","title":"AddPrompt is the string to display when asking if another value should be entered."},"disabled":{"type":"string","title":"Disabled will disable this prompt if set to true."},"name":{"type":"string","title":"Name is the name of the variable.","description":"The name can be a golang template and can use values gathered from previous\nvariables in the variables array."},"required":{"type":"boolean","title":"Required whether or not a non-empty value is required."},"variables":{"type":"array","title":"TemplateVariables holds the variables that make up the fields of the object.","items":{"type":"object","title":"TemplateVariable is the base interface for a variable.","anyOf":[{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-SimpleVar"},{"$ref":"#"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-CustomizedVar"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-Selection"},{"$ref":"#/definitions/github_com-brainicorn-skelp-skelplate-MultiValue"}]}}},"required":["name","variables"],"additionalProperties":false}`

	// GithubComBrainicornSkelpSkelplateHooks is a json-schema accessor
	GithubComBrainicornSkelpSkelplateHooks = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"Hooks is the object that holds arrays of the various hook scripts.","description":"Each lifecycle can have an array of strings that represent the shell scripts to run.\nEach string should be the basename of the script file followed by any arguments.\nThe string will be processed as a Go Template so the args can use built-in functions and any data\nthat's available from gathering input. The script is assumed to live in the template repo's hooks driectory.\n","properties":{"postGen":{"type":"array","items":{"type":"string"}},"preGen":{"type":"array","items":{"type":"string"}},"preInput":{"type":"array","items":{"type":"string"}}}}`

	// GithubComBrainicornSkelpSkelplateTemplateExclude is a json-schema accessor
	GithubComBrainicornSkelpSkelplateTemplateExclude = `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","title":"TemplateExclude holds a condition, which if true excludes the list of template paths from processing","properties":{"exclude":{"type":"string","title":"Exclude is a go template that should evalutate to true or false.","description":"If it evaluates to true, the list of template paths will be excluded from processing\n"},"paths":{"type":"array","title":"FilesOrDirs holds the paths that should be excluded when Excludes is true.","items":{"type":"string"}}}}`

)
