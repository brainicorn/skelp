package skelplate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/brainicorn/skelp/skelputil"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateDescriptor(descriptorBytes []byte) (SkelplateDescriptor, error) {
	var err error
	var skelplate SkelplateDescriptor
	var schemaValidationResult *gojsonschema.Result

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

	if err == nil {
		err = json.Unmarshal(descriptorBytes, &skelplate)
	}

	if err == nil {
		// validate hooks
		if skelplate.TemplateHooks.PreInput != nil && len(skelplate.TemplateHooks.PreInput) > 0 {
			err = validateHookArray(skelplate.TemplateHooks.PreInput)
		}

		if err == nil && skelplate.TemplateHooks.PreGen != nil && len(skelplate.TemplateHooks.PreGen) > 0 {
			err = validateHookArray(skelplate.TemplateHooks.PreGen)
		}

		if err == nil && skelplate.TemplateHooks.PostGen != nil && len(skelplate.TemplateHooks.PostGen) > 0 {
			err = validateHookArray(skelplate.TemplateHooks.PostGen)
		}
	}

	return skelplate, err
}

func validateHookArray(hooks []string) error {
	var err error
	for _, hook := range hooks {
		if len(strings.TrimSpace(hook)) < 1 {
			err = fmt.Errorf("script must not be blank", hook)
		}

		if err == nil {
			// check the first arg is a filename and is a file, not a dir
			script := strings.Split(hook, " ")[0]
			if strings.HasPrefix(script, ".") ||
				strings.HasPrefix(script, string(filepath.Separator)) ||
				strings.HasPrefix(script, "{{") ||
				strings.HasSuffix(script, string(filepath.Separator)) ||
				strings.Contains(script, ".."+string(filepath.Separator)) {
				err = fmt.Errorf("script name %s must be a base filename", script)
			}

			if err == nil {
				// make sure the template can be parsed, but don't execute it
				_, err = template.New("hookTmpl").Option(skelputil.TemplateOptions()...).Funcs(skelputil.FunctionMap()).Parse(hook)
			}
		}
	}

	if err != nil {
		err = fmt.Errorf("error parsing hooks: %s", err.Error())
	}

	return err
}
