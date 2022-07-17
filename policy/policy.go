package policy

import (
	"encoding/json"
	"log"
	"nweisenauer/dag"
)

func BuildPolicy(jsonRules []byte) ([]byte, error) {
	var rules RuleSet
	err := json.Unmarshal(jsonRules, &rules)

	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}

	sortedRuleIds, err := sortRuleIds(rules)

	if err != nil {
		log.Println(err)
		return make([]byte, 0), err
	}

	policy := convertSortedRuleSetToPolicy(sortedRuleIds, rules)
	result, err := json.MarshalIndent(policy, "", "    ")

	if err != nil {
		log.Println(err)
	}

	return result, err
}

func sortRuleIds(rules RuleSet) ([]string, error) {
	graph := convertRuleSetToGraph(rules)
	result, err := dag.TopologicalSort(graph)
	sortedRuleIds := make([]string, 0)

	if err != nil {
		log.Println(err)
		return sortedRuleIds, err
	}

	for _, ruleIdSlice := range result {
		sortedRuleIds = append(sortedRuleIds, ruleIdSlice[0])
	}

	return sortedRuleIds, nil
}

func convertRuleSetToGraph(rules RuleSet) dag.Graph {
	graph := make(dag.Graph, len(rules.Rules))

	for _, rule := range rules.Rules {
		graph[rule.Id] = make([]string, 0)
		for _, dependency := range rule.Requires {
			graph[rule.Id] = append(graph[rule.Id], dependency)
		}
	}

	return graph
}

func convertSortedRuleSetToPolicy(sortedRuleIds []string, rules RuleSet) Policy {
	policy := make(Policy, len(sortedRuleIds))

	for i, ruleId := range sortedRuleIds {
		for _, rule := range rules.Rules {
			if rule.Id == ruleId {
				policy[i] = PolicyRule{
					Id:   rule.Id,
					Head: rule.Head,
					Body: rule.Body,
				}
			}
		}
	}

	return policy
}
