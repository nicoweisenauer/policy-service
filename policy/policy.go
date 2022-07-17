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

	sortedRuleIds, err := sortRuleIds(rules)

	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}

	responseRules := convertToResponseRuleFormat(sortedRuleIds, rules)
	result, err := json.MarshalIndent(responseRules, "", "    ")

	if err != nil {
		log.Println(err)
	}

	return result, err
}

func sortRuleIds(rules RuleSet) ([]string, error) {
	graph := convertRuleSetToGraphRepresentation(rules)
	result, err := dag.TopologicalSort(graph)
	sortedRuleIds := make([]string, 0)

	if err != nil {
		log.Println(err)
		return sortedRuleIds, err
	}

	fmt.Println(result)

	for _, ruleIdSlice := range result {
		sortedRuleIds = append(sortedRuleIds, ruleIdSlice[0])
	}

	return sortedRuleIds, nil
}

func convertRuleSetToGraphRepresentation(rules RuleSet) map[string][]string {
	graph := make(map[string][]string, len(rules.Rules))

	for _, rule := range rules.Rules {
		graph[rule.Id] = make([]string, 0)
		for _, dependency := range rule.Requires {
			graph[rule.Id] = append(graph[rule.Id], dependency)
		}
	}

	return graph
}

func convertToResponseRuleFormat(sortedRuleIds []string, rules RuleSet) []ResultRule {
	responseRules := make([]ResultRule, len(sortedRuleIds))

	for i, ruleId := range sortedRuleIds {
		for _, rule := range rules.Rules {
			if rule.Id == ruleId {
				responseRules[i] = ResultRule{
					Id:   rule.Id,
					Head: rule.Head,
					Body: rule.Body,
				}
			}
		}
	}

	return responseRules
}
