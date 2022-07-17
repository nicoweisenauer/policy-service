package policy

import (
	"encoding/json"
	"fmt"
	"log"
	"nweisenauer/dag"
)

func ProcessPolicy(policy []byte) ([]byte, error) {
	var rules RuleSet
	err := json.Unmarshal(policy, &rules)

	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}

	rules, err = sortRules(rules)

	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}

	responseRules := make([]Rule, len(rules.Rules))

	for i, r := range rules.Rules {
		responseRules[i] = Rule{
			Id:   r.Id,
			Head: r.Head,
			Body: r.Body,
		}
	}

	result, err := json.MarshalIndent(responseRules, "", "    ")

	if err != nil {
		log.Println(err)
	}

	return result, err
}

func sortRules(rules RuleSet) (RuleSet, error) {
	graph := make(map[string][]string, len(rules.Rules))

	for _, rule := range rules.Rules {
		graph[rule.Id] = make([]string, 0)
		for _, dependency := range rule.Requires {
			graph[rule.Id] = append(graph[rule.Id], dependency)
		}
	}

	result, err := dag.TopologicalSort(graph)
	resultRules := RuleSet{}

	if err != nil {
		log.Println(err)
		return resultRules, err
	}

	fmt.Println(result)

	for _, ruleIdSlice := range result {
		for _, rule := range rules.Rules {
			if rule.Id == ruleIdSlice[0] {
				resultRules.Rules = append(resultRules.Rules, rule)
			}
		}
	}

	return resultRules, nil
}
