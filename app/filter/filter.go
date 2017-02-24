// Package filter defines functions that check if a kill is within a given channel's
// configured filters.
package filter

import (
	"github.com/eveopsec/zk2s/app/config"
	"github.com/vivace-io/evelib/crest"
	"github.com/vivace-io/evelib/redisq"
)

// Within returns true if the kill is within the channel's filter, false otherwise
func Within(payload redisq.Payload, channel config.Channel) bool {
	// KillID of 0 should be skipped as it means no kill was returned from RedisQ
	if payload.KillID == 0 {
		return false
	}
	if !valueOK(payload, channel) {
		return false
	}
	if !shipOK(payload, channel) {
		return false
	}
	if !involvedOK(payload, channel) {
		return false
	}
	return true
}

// IsLoss returns true if the kill is a loss, false otherwise
func IsLoss(payload redisq.Payload, channel config.Channel) bool {
	if characterOK(payload.Killmail.Victim.Character, channel) {
		return true
	}
	if corporationOK(payload.Killmail.Victim.Corporation, channel) {
		return true
	}
	if allianceOK(payload.Killmail.Victim.Alliance, channel) {
		return true
	}
	return false
}

// IsAwox returns true if the kill was an Awox, partial or otherwise. False if not.
func IsAwox(payload redisq.Payload, channel config.Channel) bool {
	if IsLoss(payload, channel) {
		for a := range payload.Killmail.Attackers {
			if characterOK(payload.Killmail.Attackers[a].Character, channel) {
				return true
			}
			if corporationOK(payload.Killmail.Attackers[a].Corporation, channel) {
				return true
			}
			if allianceOK(payload.Killmail.Attackers[a].Alliance, channel) {
				return true
			}
		}
	}
	return false
}

func valueOK(payload redisq.Payload, channel config.Channel) bool {
	// If kill value is within [MinimumValue, MaximumValue] return true
	if payload.Zkb.TotalValue >= float32(channel.MinimumValue) && payload.Zkb.TotalValue <= float32(channel.MaximumValue) {
		return true
	}
	// If kill value is greater than min AND no max value is set, return true
	if payload.Zkb.TotalValue >= float32(channel.MinimumValue) && channel.MaximumValue <= 0 {
		return true
	}
	return false
}

func shipOK(payload redisq.Payload, channel config.Channel) bool {
	for ship := range channel.ExcludedShips {
		if payload.Killmail.Victim.ShipType.Name == channel.ExcludedShips[ship] {
			return false
		}
		if string(payload.Killmail.Victim.ShipType.ID) == channel.ExcludedShips[ship] {
			return false
		}
	}
	return true
}

func involvedOK(payload redisq.Payload, channel config.Channel) bool {
	if characterOK(payload.Killmail.Victim.Character, channel) {
		return true
	}
	if corporationOK(payload.Killmail.Victim.Corporation, channel) {
		return true
	}
	if allianceOK(payload.Killmail.Victim.Alliance, channel) {
		return true
	}
	for a := range payload.Killmail.Attackers {
		if characterOK(payload.Killmail.Attackers[a].Character, channel) {
			return true
		}
		if corporationOK(payload.Killmail.Attackers[a].Corporation, channel) {
			return true
		}
		if allianceOK(payload.Killmail.Attackers[a].Alliance, channel) {
			return true
		}
	}
	return false
}

func characterOK(char crest.Character, channel config.Channel) bool {
	for c := range channel.IncludeCharacters {
		if (char.Name == channel.IncludeCharacters[c]) || (string(char.ID) == channel.IncludeCharacters[c]) {
			return true
		}
	}
	return false
}

func corporationOK(corp crest.Corporation, channel config.Channel) bool {
	for c := range channel.IncludeCorporations {
		if corp.Name == channel.IncludeCorporations[c] {
			return true
		}
	}
	return false
}

func allianceOK(alli crest.Alliance, channel config.Channel) bool {
	for a := range channel.IncludeAlliances {
		if alli.Name == channel.IncludeAlliances[a] {
			return true
		}
	}
	return false
}
