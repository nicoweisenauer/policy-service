package main

import (
	"fmt"
	"io"
	"net/http"
	"nweisenauer/policy"
)

func main() {
	http.HandleFunc("/", policyService)
	http.ListenAndServe(":8080", nil)
}

func policyService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: policyService")
	requestBody, _ := io.ReadAll(r.Body)
	returnPolicy, _ := policy.ProcessPolicy(requestBody)
	w.Write(returnPolicy)
}
