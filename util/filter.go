package util

import (
	"github.com/eveopsec/zk2s/config"
	"github.com/vivace-io/evelib/crest"
	"github.com/vivace-io/evelib/zkill"
)

/* util/filter.go
 * Defines functions for filtering through kills.
 */

// WithinFilter returns true if the kill is within the channel's filter, false otherwise
func WithinFilter(kill *zkill.Kill, channel config.Channel) bool {
	// KillID of 0 should be skipped as it means no kill was returned from RedisQ
	if kill.KillID == 0 {
		return false
	}
	if !valueOK(kill, channel) {
		return false
	}
	if !shipOK(kill, channel) {
		return false
	}
	if !involvedOK(kill, channel) {
		return false
	}
	return true
}

// IsLoss returns true if the kill is a loss, false otherwise
func IsLoss(kill *zkill.Kill, channel config.Channel) bool {
	if characterOK(kill.Killmail.Victim.Character, channel) {
		return true
	}
	if corporationOK(kill.Killmail.Victim.Corporation, channel) {
		return true
	}
	if allianceOK(kill.Killmail.Victim.Alliance, channel) {
		return true
	}
	return false
}

// IsAwox returns true if the kill was an Awox, partial or otherwise. False if not.
func IsAwox(kill *zkill.Kill, channel config.Channel) bool {
	if IsLoss(kill, channel) {
		for a := range kill.Killmail.Attackers {
			if characterOK(kill.Killmail.Attackers[a].Character, channel) {
				return true
			}
			if corporationOK(kill.Killmail.Attackers[a].Corporation, channel) {
				return true
			}
			if allianceOK(kill.Killmail.Attackers[a].Alliance, channel) {
				return true
			}
		}
	}
	return false
}

func valueOK(kill *zkill.Kill, channel config.Channel) bool {
	// If kill value is within [MinimumValue, MaximumValue] return true
	if kill.Zkb.TotalValue >= float32(channel.MinimumValue) && kill.Zkb.TotalValue <= float32(channel.MaximumValue) {
		return true
	}
	// If kill value is greater than min AND no max value is set, return true
	if kill.Zkb.TotalValue >= float32(channel.MinimumValue) && channel.MaximumValue == 0 {
		return true
	}
	return false
}

// returns true if ship is NOT excluded by name, false otherwise
func shipOK(kill *zkill.Kill, channel config.Channel) bool {
	for ship := range channel.ExcludedShips {
		if kill.Killmail.Victim.ShipType.Name == channel.ExcludedShips[ship] {
			return false
		}
		if string(kill.Killmail.Victim.ShipType.ID) == channel.ExcludedShips[ship] {
			return false
		}
	}
	return true
}

func involvedOK(kill *zkill.Kill, channel config.Channel) bool {
	if characterOK(kill.Killmail.Victim.Character, channel) {
		return true
	}
	if corporationOK(kill.Killmail.Victim.Corporation, channel) {
		return true
	}
	if allianceOK(kill.Killmail.Victim.Alliance, channel) {
		return true
	}
	for a := range kill.Killmail.Attackers {
		if characterOK(kill.Killmail.Attackers[a].Character, channel) {
			return true
		}
		if corporationOK(kill.Killmail.Attackers[a].Corporation, channel) {
			return true
		}
		if allianceOK(kill.Killmail.Attackers[a].Alliance, channel) {
			return true
		}
	}
	return false
}

// returns true if passed character is in channel.IncludeCharacters, false otherwise
func characterOK(char crest.Character, channel config.Channel) bool {
	for c := range channel.IncludeCharacters {
		if (char.Name == channel.IncludeCharacters[c]) || (string(char.ID) == channel.IncludeCharacters[c]) {
			return true
		}
	}
	return false
}

// returns true if passed corporation is in channel.IncludeCorporations, false otherwise
func corporationOK(corp crest.Corporation, channel config.Channel) bool {
	for c := range channel.IncludeCorporations {
		if corp.Name == channel.IncludeCorporations[c] {
			return true
		}
	}
	return false
}

// returns true if passed alliance is in channel.IncludeAlliances, false otherwise
func allianceOK(alli crest.Alliance, channel config.Channel) bool {
	for a := range channel.IncludeAlliances {
		if alli.Name == channel.IncludeAlliances[a] {
			return true
		}
	}
	return false
}
