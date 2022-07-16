package main

type RuleSet struct {
	Rules []struct {
		Id       string   `json:"id"`
		Head     string   `json:"head"`
		Body     string   `json:"body"`
		Requires []string `json:"requires"`
	} `json:"rules"`
}
