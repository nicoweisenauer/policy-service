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
	result, err := ProcessPolicy(request)
	if !json.Valid(result) || err != nil {
		t.Fatalf(`TestMainReturnsValidJSON: %q, %v, no valid JSON`, result, err)
	}
}

func TestProcessPolicyReturnsAllRules(t *testing.T) {
	result, err := ProcessPolicy(request)
	if err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, no valid JSON`, result, err)
	}

	var rules []Rule
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

	result, err := ProcessPolicy(request)
	if err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, no valid JSON`, result, err)
	}

	var expectedRules []Rule
	err = json.Unmarshal(expected, &expectedRules)

	var resultRules []Rule
	err = json.Unmarshal(result, &resultRules)
	if (len(resultRules) != 4) || err != nil {
		t.Fatalf(`TestProcessPolicySortsRules: result length is not correct`)
	}
	if !(expectedRules[0] == resultRules[0]) {
		t.Fatalf(`TestProcessPolicySortsRules: %q, want match for %#q`, expectedRules[0], resultRules[0])
	}
	if !(expectedRules[1] == resultRules[1]) {
		t.Fatalf(`TestProcessPolicySortsRules: %q, want match for %#q`, expectedRules[1], resultRules[1])
	}
	if !(expectedRules[2] == resultRules[2]) {
		t.Fatalf(`TestProcessPolicySortsRules: %q, want match for %#q`, expectedRules[2], resultRules[2])
	}
	if !(expectedRules[3] == resultRules[3]) {
		t.Fatalf(`TestProcessPolicySortsRules: %q, want match for %#q`, expectedRules[3], resultRules[4])
	}
}
