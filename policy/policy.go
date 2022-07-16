package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	result, error := json.Marshal(rules)

	if error != nil {
		log.Fatal(error)
	}

	return result, error
}
