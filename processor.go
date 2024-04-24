package main

import (
	"time"
)

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

func getLatestSuccessTimestamp(proofs []SuccintProof) time.Time {
	var latestSuccessTime time.Time

	for _, proof := range proofs {
		parsedTime, err := time.Parse(time.RFC3339, proof.CreatedAt)
		if err != nil {
			continue
		}

		if proof.Status == "success" && parsedTime.After(latestSuccessTime) {
			latestSuccessTime = parsedTime
		}
	}
	return latestSuccessTime
}

func getLatestFailureTimestamp(proofs []SuccintProof) time.Time {
	var latestFailureTime time.Time

	for _, proof := range proofs {
		parsedTime, err := time.Parse(time.RFC3339, proof.CreatedAt)
		if err != nil {
			continue
		}

		if proof.Status == "failure" && parsedTime.After(latestFailureTime) {
			latestFailureTime = parsedTime
		}
	}
	return latestFailureTime
}
