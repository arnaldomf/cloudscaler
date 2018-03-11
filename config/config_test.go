package config

import (
	"testing"
)

var data = `
{
	"EventName": "blackfriday",
	"AutoScalingGroups": [
		{
			"Name": "backend-asg",
			"Growth": 30,
			"Region": "sa-east-1",
			"Profile": "default"
		},
		{
			"Name": "frontend-asg",
			"Growth": 10,
			"Region": "sa-east-1"
		}
	]
}
`

var validDataWithoutEventName = `
{
	"AutoScalingGroups": [
		{
			"Name": "backend-asg",
			"Growth": 30,
			"Region": "sa-east-1",
			"Profile": "default"
		},
		{
			"Name": "frontend-asg",
			"Growth": 10,
			"Region": "sa-east-1"
		}
	]
}`

var invalidData = `
{
	"EventName": "xpto"
	"AutoScalingGroups": []
}
`

func TestNewConfiguration(t *testing.T) {
	configuration, err := NewConfiguration([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if configuration.EventName != "blackfriday" {
		t.Errorf("Expected EventName blackfriday, received '%s'",
			configuration.EventName)
	}
	if len(configuration.AutoScalingGroups) != 2 {
		t.Errorf("Expected 2 autoscaling grroups, received %d",
			len(configuration.AutoScalingGroups))
	}
	if configuration.AutoScalingGroups[0].Name != "backend-asg" ||
		configuration.AutoScalingGroups[1].Name != "frontend-asg" {
		t.Error("Wrong data on autoscaling groups")
	}
}

func TestNewConfigurationWithInvalidJSON(t *testing.T) {
	_, err := NewConfiguration([]byte(invalidData))
	if err == nil {
		t.Error("Invalid JSON must return error")
	}
}

func TestNewConfigurationWithoutRequiredField(t *testing.T) {
	_, err := NewConfiguration([]byte(validDataWithoutEventName))
	if err == nil {
		t.Error("Missing field, must return error")
	}
}

func TestInvalidConfiguration(t *testing.T) {
	configuration := new(Configuration)
	if err := validate(configuration); err == nil {
		t.Error("Empty configuration, should have returned an error")
	}
	configuration.EventName = "Blackfriday"
	if err := validate(configuration); err == nil {
		t.Error("Empty AutoScalingGroups, should have returned an error")
	}

	autoScalingGroups := make([]*AutoScalingGroup, 0)

	autoScalingGroups = append(autoScalingGroups, &AutoScalingGroup{
		Name:    "backend-asg",
		Growth:  30,
		Region:  "sa-east-1",
		Profile: "default",
	})
	configuration.AutoScalingGroups = autoScalingGroups
	if err := validate(configuration); err != nil {
		t.Fatalf("Valid autoscaling group returned: %v", err)
	}

	autoScalingGroups = append(autoScalingGroups, &AutoScalingGroup{
		Growth:  30,
		Region:  "sa-east-1",
		Profile: "default",
	})
	configuration.AutoScalingGroups = autoScalingGroups
	if err := validate(configuration); err == nil {
		t.Error("AutoScaling without name should have returned error")
	}
	configuration.AutoScalingGroups[len(configuration.AutoScalingGroups)-1].Name =
		"group"
	configuration.AutoScalingGroups[len(configuration.AutoScalingGroups)-1].Growth =
		0
	if err := validate(configuration); err == nil {
		t.Error("AutoScaling with 0 growth should have returned error")
	}
	configuration.AutoScalingGroups[len(configuration.AutoScalingGroups)-1].Growth =
		-1
	if err := validate(configuration); err == nil {
		t.Error("AutoScaling with negative growth should have returned error")
	}
}
