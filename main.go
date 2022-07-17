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
	returnPolicy, err := policy.ProcessPolicy(requestBody)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(returnPolicy)
}
