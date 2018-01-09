package provider

// ReplayProvider is an interface that provides read/write methods for dealing with replay data.
type ReplayProvider interface {
	WriteData(data map[string]interface{}, projectRoot, templatePath string) (string, error)
	ReadData(projectRoot, templatePath string) (map[string]interface{}, error)
}

type DefaultReplayProvider struct{}

func (drp *DefaultReplayProvider) WriteData(data map[string]interface{}, projectRoot, templatePath string) (string, error) {
	return "", nil
}

func (drp *DefaultReplayProvider) ReadData(projectRoot, templatePath string) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}
