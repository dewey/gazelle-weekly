package applemusic

// releaseYearWithinRange checks if a release year is within a reasonable range, this helps us filter out similarly named
// albums on streaming services
func releaseYearWithinRange(releaseYearCandidate int, releaseYearInput, rangeYears int) bool {
	if releaseYearCandidate >= releaseYearInput-rangeYears && releaseYearCandidate <= releaseYearInput+rangeYears {
		return true
	}
	return false
}
