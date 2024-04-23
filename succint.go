package main

type SuccintProof struct {
	ID           string      `json:"id"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
	ProofRequest interface{} `json:"proof_request"`
	ProofRelease interface{} `json:"proof_release"`
	Edges        interface{} `json:"edges"`
	Requests     interface{} `json:"requests"`
}
