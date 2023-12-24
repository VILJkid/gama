package yaml

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type WorkflowContent struct {
	Name string `yaml:"name"`
	On   struct {
		WorkflowDispatch struct {
			Inputs map[string]WorkflowInput `yaml:"inputs"`
		} `yaml:"workflow_dispatch"`
	} `yaml:"on"`
}

type WorkflowInput struct {
	Description string            `yaml:"description"`
	Required    bool              `yaml:"required"`
	Default     interface{}       `yaml:"default,omitempty"`
	Type        string            `yaml:"type,omitempty"`
	Options     []string          `yaml:"options,omitempty"`
	JSONContent map[string]string `yaml:"-"` // This field is for internal use and won't be filled directly by the YAML unmarshaler
}

func (i *WorkflowInput) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Define a shadow type to avoid recursion
	type shadow WorkflowInput
	if err := unmarshal((*shadow)(i)); err != nil {
		return err
	}

	// Attempt to unmarshal JSON content if the default value is a string
	if defaultStr, ok := i.Default.(string); ok {
		tempMap := make(map[string]string)
		if err := json.Unmarshal([]byte(defaultStr), &tempMap); err == nil {
			i.JSONContent = tempMap
		}
	}
	return nil
}

func UnmarshalWorkflowContent(data []byte) (*WorkflowContent, error) {
	var workflow WorkflowContent
	err := yaml.Unmarshal(data, &workflow)
	if err != nil {
		return nil, err
	}

	return &workflow, nil
}