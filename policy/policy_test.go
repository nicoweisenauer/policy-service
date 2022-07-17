package policy

import (
	"encoding/json"
	"testing"
)

var request = []byte(`
{
	"rules": [
		{
			"id": "rule-1",
			"head": "default allow = false"
		},
		{
			"id": "rule-2",
			"head": "allow",
			"body": "method == \"GET\"; data.roles[\"dev\"][_] == input.user",
			"requires": [
				"rule-3",
				"rule-4"
			]
		},
		{
			"id": "rule-3",
			"head": "allow",
			"body": "input.user == \"alice\"",
			"requires": [
				"rule-1"
			]
		},
		{
			"id": "rule-4",
			"head": "allow",
			"body": "input.user == \"bob\"; method == \"GET\"",
			"requires": [
				"rule-3"
			]	
		}
	]
}`)

func TestProcessPolicyReturnsValidJSON(t *testing.T) {
	result, err := BuildPolicy(request)
	if !json.Valid(result) || err != nil {
		t.Fatalf(`TestMainReturnsValidJSON: %q, %v, no valid JSON`, result, err)
	}
}

func TestProcessPolicyReturnsAllRules(t *testing.T) {
	result, err := BuildPolicy(request)
	if err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, no valid JSON`, result, err)
	}

	var rules Policy
	err = json.Unmarshal(result, &rules)

	for _, id := range []string{"rule-1", "rule-2", "rule-3", "rule-4"} {
		contained := false
		for _, rule := range rules {
			if rule.Id == id {
				contained = true
				break
			}
		}
		if !contained {
			t.Fatalf(`TestMainReturnsAllRules: %q, %v, rule is missing from result`, result, err)
		}
	}

	if len(rules) != 4 || err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, result rules do not match requested rules`, result, err)
	}
}

func TestProcessPolicySortsRules(t *testing.T) {
	expected := []byte(`
	[
		{
			"id": "rule-1",
			"head": "default allow = false",
			"body": ""
		},
		{
			"id": "rule-3",
			"head": "allow",
			"body": "input.user == \"alice\""
		},
		{
			"id": "rule-4",
			"head": "allow",
			"body": "input.user == \"bob\"; method == \"GET\""
		},
		{
			"id": "rule-2",
			"head": "allow",
			"body": "method == \"GET\"; data.roles[\"dev\"][_] == input.user"
		}
	]`)

	result, err := BuildPolicy(request)
	if err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, no valid JSON`, result, err)
	}

	var expectedRules Policy
	err = json.Unmarshal(expected, &expectedRules)
	expectedRulesJSON, _ := json.MarshalIndent(expectedRules, "", "    ")

	if (string(result) != string(expectedRulesJSON)) || err != nil {
		t.Fatalf(`TestProcessPolicySortsRules: result json does not match expected json, result: %q`, string(result))
	}
}

func TestProcessPolicyDetectsCyclicDependencies(t *testing.T) {
	cyclicRequest := []byte(`
	{
		"rules": [
			{
				"id": "rule-1",
				"requires": [
					"rule-2"
				]
			},
			{
				"id": "rule-2",
				"requires": [
					"rule-1"
				]
			}
		]
	}`)

	result, err := BuildPolicy(cyclicRequest)

	if len(result) != 0 {
		t.Fatalf(`TestProcessPolicySortsRules: result should be empty, but is %q`, string(result))
	}
	if !(err.Error() == "Cyclic dependency detected!") {
		t.Fatalf(`TestProcessPolicySortsRules: should return error about cyclic dependency`)
	}

}
