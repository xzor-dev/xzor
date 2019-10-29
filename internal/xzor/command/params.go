package command

import "fmt"

// Params is an arbitrary map of parameters used when executing commands.
type Params map[string]interface{}

// String returns a parameter as a string.
func (p Params) String(name string) (string, error) {
	if p[name] == nil {
		return "", fmt.Errorf("invalid parameter name: %s", name)
	}
	str, ok := p[name].(string)
	if !ok {
		return "", fmt.Errorf("could not convert parameter '%s' to a string", name)
	}
	return str, nil
}
