package aino

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

func (c *char) a1BurstEnhance() (int, attacks.ICDGroup) {
	if c.Base.Ascension < 1 {
		return 90, attacks.ICDGroupDefault
	}

	if c.getMoonsignLevel() < 2 {
		return 90, attacks.ICDGroupDefault
	}
	return 42, attacks.ICDGroupAinoBurstEnhanced
}

func (c *char) a4Dmg() float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	return c.Stat(attributes.EM) * 0.5
}
