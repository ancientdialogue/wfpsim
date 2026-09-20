package hacks

import "github.com/genshinsim/gcsim/pkg/core"

const (
	stellarSwirlDuration    = 8 * 60
	StellarSwirlBonusDurKey = "stellar-swirl-bonus-dur"
)

func StellarSwirlDur(c *core.Core) int {
	return stellarSwirlDuration + int(c.Flags.Custom[StellarSwirlBonusDurKey])
}
