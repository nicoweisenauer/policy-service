package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nweisenauer/dag"
)

func main() {

	request := []byte(`
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

	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("ProcessPolicy: ")
	log.SetFlags(0)

	result, err := ProcessPolicy(request)

	if err != nil {
		log.Fatal(err)
	}

	stringResult := string(result[:])
	// If no error was returned, print the result to the console.
	fmt.Println(stringResult)
}

func ProcessPolicy(policy []byte) ([]byte, error) {
	var rules RuleSet
	error := json.Unmarshal(policy, &rules)

	if error != nil {
		log.Fatal(error)
	}

	rules = sortRules(rules)

	responseRules := make([]Rule, len(rules.Rules))

	for i, r := range rules.Rules {
		responseRules[i] = Rule{
			Id:   r.Id,
			Head: r.Head,
			Body: r.Body,
		}
	}

	result, error := json.Marshal(responseRules)

	if error != nil {
		log.Fatal(error)
	}

	return result, error
}

func sortRules(rules RuleSet) RuleSet {
	graph := make(map[string][]string, len(rules.Rules))

	for _, rule := range rules.Rules {
		graph[rule.Id] = make([]string, 0)
		for _, dependency := range rule.Requires {
			graph[rule.Id] = append(graph[rule.Id], dependency)
		}
	}

	result := dag.Connections(graph)

	fmt.Println(result)

	resultRules := RuleSet{}

	for _, ruleIdSlice := range result {
		for _, rule := range rules.Rules {
			if rule.Id == ruleIdSlice[0] {
				resultRules.Rules = append(resultRules.Rules, rule)
			}
		}
	}

	return resultRules
}
