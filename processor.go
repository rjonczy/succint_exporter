package main

func countFailedProofs(proofs []SuccintProof) int {
	count := 0
	for _, proof := range proofs {
		if proof.Status == "failure" {
			count++
		}
	}
	return count
}

func countRunningProofs(proofs []SuccintProof) int {
	count := 0
	for _, proof := range proofs {
		if proof.Status == "running" {
			count++
		}
	}
	return count
}
