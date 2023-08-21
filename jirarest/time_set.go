package jirarest

import "strings"

// TimeRemaining returns calculated timeRemainingOriginal and timeRemaiing and against the timeOriginalEstimate and timeEstimate respectively.
func TimeRemaining(status string, timeOriginalEstimate, timeEstimate, timeSpent int) (timeRemainingOriginal, timeRemaining int) {
	if timeOriginalEstimate < 0 {
		timeOriginalEstimate = 0
	}
	if timeEstimate < 0 {
		timeEstimate = 0
	}
	if timeSpent < 0 {
		timeEstimate = 0
	}
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "closed" || status == "done" {
		return 0, 0
	}
	timeRemaining = timeEstimate - timeSpent
	if timeRemaining < 0 {
		timeRemaining = 0
	}
	timeRemainingOriginal = timeOriginalEstimate - timeSpent
	if timeRemainingOriginal < 0 {
		timeRemainingOriginal = 0
	}
	return
}
