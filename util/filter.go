package util

import (
	"strconv"

	"github.com/spf13/viper"
	"github.com/vivace-io/evelib/crest"
	"github.com/vivace-io/evelib/zkill"
)

/* util/filter.go
 * Defines functions for filtering through kills.
 */

// isWithinFilters returns true if the kill is within the defined filters.
func isWithinFilters(kill *zkill.ZKill, config *viper.Viper) bool {
	if kill.Zkb.TotalValue < float32(config.GetInt("iskMinimum")) {
		return false
	}
	if !withinShipFilter(kill.Killmail.Victim.ShipType, config) {
		return false
	}
	if withinCorpFilter(kill.Killmail.Victim.Corporation.ID, config) {
		if config.GetBool("excludeLosses") {
			return false
		}
		return true
	}
	if withinAllianceFilter(kill.Killmail.Victim.Alliance.ID, config) {
		if config.GetBool("excludeLosses") {
			return false
		}
		return true
	}
	for x := range kill.Killmail.Attackers {
		if withinCorpFilter(kill.Killmail.Attackers[x].Corporation.ID, config) {
			return true
		}
		if withinAllianceFilter(kill.Killmail.Attackers[x].Corporation.ID, config) {
			return true
		}
	}
	return false
}

func withinShipFilter(ship crest.Type, config *viper.Viper) bool {
	ships := config.GetStringSlice("excludeShips")
	if len(ships) == 0 {
		return true
	}
	for x := range ships {
		if strconv.Itoa(ship.ID) == ships[x] {
			return false
		}
		if ship.Name == ships[x] {
			return false
		}
	}
	return true
}

func withinCorpFilter(id int, config *viper.Viper) bool {
	corps := config.GetStringSlice("includeCorps")
	if len(corps) == 0 {
		return false
	}
	for x := range corps {
		if strconv.Itoa(id) == corps[x] {
			return true
		}
	}
	return false
}

func withinAllianceFilter(id int, config *viper.Viper) bool {
	alliances := config.GetStringSlice("includeAlliances")
	if len(alliances) == 0 {
		return false
	}
	for x := range alliances {
		if strconv.Itoa(id) == alliances[x] {
			return true
		}
	}
	return false
}
