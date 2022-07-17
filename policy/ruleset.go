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

type ResultRule struct {
	Id   string `json:"id"`
	Head string `json:"head"`
	Body string `json:"body"`
}
