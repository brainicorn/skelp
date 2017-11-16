package provider

// DataProvider is a function that returns the data to be applied to a template or an error
type DataProvider func(templateRoot string) (interface{}, error)
