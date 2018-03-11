package config

import (
	"encoding/json"
	"fmt"
)

// AutoScalingGroup defines an ASG that needs to grow before an event
type AutoScalingGroup struct {
	Name    string `json:"Name"`
	Growth  int    `json:"Growth"`
	Region  string `json:"Region"`
	Profile string `json:"Profile,omitempty"`
}

// Configuration defines the data coming from configurantion file
type Configuration struct {
	EventName         string              `json:"EventName"`
	AutoScalingGroups []*AutoScalingGroup `json:"AutoScalingGroups"`
}

// NewConfiguration unmarshals data and validates configuration fields
func NewConfiguration(data []byte) (*Configuration, error) {
	tag := "NewConfiguration:"
	configuration := new(Configuration)
	if err := json.Unmarshal(data, configuration); err != nil {
		return nil, fmt.Errorf("%s %v", tag, err)
	}

	if err := validate(configuration); err != nil {
		return nil, fmt.Errorf("%s %v", tag, err)
	}
	return configuration, nil
}

func validate(configuration *Configuration) error {
	tag := "validateConfiguration:"
	if len(configuration.EventName) == 0 {
		return fmt.Errorf("%s EventName can't be empty", tag)
	}
	if len(configuration.AutoScalingGroups) == 0 {
		return fmt.Errorf("%s AutoScalingGroups can't be empty", tag)
	}
	for _, autoScalingGroup := range configuration.AutoScalingGroups {
		if err := validateAutoScalingGroup(autoScalingGroup); err != nil {
			return fmt.Errorf("%s %v", tag, err)
		}
	}
	return nil
}

func validateAutoScalingGroup(autoScalingGroup *AutoScalingGroup) error {
	tag := "validateAutoScalingGroup:"
	if len(autoScalingGroup.Name) == 0 {
		return fmt.Errorf("%s Name can't be empty", tag)
	}
	if autoScalingGroup.Growth <= 0 {
		return fmt.Errorf("%s Growth must be greater than 0", tag)
	}
	return nil
}
