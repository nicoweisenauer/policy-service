package main

import (
	"fmt"
	"io"
	"net/http"
	"nweisenauer/policy"
)

// 	request := []byte(`
// 	{
// 		"rules": [
// 			{
// 				"id": "rule-1",
// 				"head": "default allow = false"
// 			},
// 			{
// 				"id": "rule-2",
// 				"head": "allow",
// 				"body": "method == \"GET\"; data.roles[\"dev\"][_] == input.user",
// 				"requires": [
// 					"rule-3",
// 					"rule-4"
// 				]
// 			},
// 			{
// 				"id": "rule-3",
// 				"head": "allow",
// 				"body": "input.user == \"alice\"",
// 				"requires": [
// 					"rule-1"
// 				]
// 			},
// 			{
// 				"id": "rule-4",
// 				"head": "allow",
// 				"body": "input.user == \"bob\"; method == \"GET\"",
// 				"requires": [
// 					"rule-3"
// 				]
// 			}
// 		]
// 	}`)

func main() {
	http.HandleFunc("/", policyService)
	http.ListenAndServe(":8080", nil)
}

func policyService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: policyService")
	requestBody, _ := io.ReadAll(r.Body)
	returnPolicy, err := policy.ProcessPolicy(requestBody)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(returnPolicy)
}
