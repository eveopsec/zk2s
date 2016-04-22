package util

import "github.com/vivace-io/evelib/zkill"

/* util/filter.go
 * Defines functions for filtering through kills.
 */

// WithinFilter returns true if the kill is within the channel's filter, false otherwise
func WithinFilter(kill *zkill.ZKill, channel *Channel) bool {
	if !valueOK(kill, channel) {
		return false
	}
	if !shipOK(kill, channel) {
		return false
	}
	if !characterOK(kill, channel) {
		return false
	}
	if !corporationOK(kill, channel) {
		return false
	}
	if !allianceOK(kill, channel) {
		return false
	}
	return true
}

// IsLoss returns true if the kill is a loss, false otherwise
func IsLoss(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

// IsAwox returns true if the kill was an Awox, partial or otherwise. False if not.
func IsAwox(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

func valueOK(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

func shipOK(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

func characterOK(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

func corporationOK(kill *zkill.ZKill, channel *Channel) bool {
	return false
}

func allianceOK(kill *zkill.ZKill, channel *Channel) bool {
	return false
}
