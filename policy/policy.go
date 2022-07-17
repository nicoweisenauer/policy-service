package policy

import (
	"encoding/json"
	"fmt"
	"log"
	"nweisenauer/dag"
)

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
