package main

import (
	"io"
	"net/http"
	"nweisenauer/policy"
)

func main() {
	http.HandleFunc("/", policyService)
	http.ListenAndServe(":8080", nil)
}

func policyService(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := io.ReadAll(r.Body)
	returnPolicy, err := policy.BuildPolicy(requestBody)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(returnPolicy)
}
