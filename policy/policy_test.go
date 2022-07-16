package main

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

func TestMainReturnsValidJSON(t *testing.T) {
	result, err := ProcessPolicy(request)
	if !json.Valid(result) || err != nil {
		t.Fatalf(`TestMainReturnsValidJSON: %q, %v, no valid JSON`, result, err)
	}
}

func TestMainReturnsAllRules(t *testing.T) {
	result, err := ProcessPolicy(request)
	if err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, no valid JSON`, result, err)
	}

	var rules RuleSet
	err = json.Unmarshal(result, &rules)

	for _, id := range []string{"rule-1", "rule-2", "rule-3", "rule-4"} {
		contained := false
		for _, rule := range rules.Rules {
			if rule.Id == id {
				contained = true
				break
			}
		}
		if !contained {
			t.Fatalf(`TestMainReturnsAllRules: %q, %v, rule is missing from result`, result, err)
		}
	}

	if len(rules.Rules) != 4 || err != nil {
		t.Fatalf(`TestMainReturnsAllRules: %q, %v, result rules do not match requested rules`, result, err)
	}
}
