package policy

type RuleSet struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Id       string   `json:"id"`
	Head     string   `json:"head"`
	Body     string   `json:"body"`
	Requires []string `json:"requires"`
}

type Policy []PolicyRule

type PolicyRule struct {
	Id   string `json:"id,omitempty"`
	Head string `json:"head,omitempty"`
	Body string `json:"body,omitempty"`
}
