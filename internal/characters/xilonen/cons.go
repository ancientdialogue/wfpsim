package xilonen

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2key = "xilonen-c2"

func (c *char) c1DurMod() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.45
}

func (c *char) c1IntervalMod() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 1.5 // make it round to 18 instead of 17.4 -> 17
}

var c2BuffGeo []float64
var c2BuffPyro []float64
var c2BuffHydro []float64
var c2BuffCryo []float64

func c2buffsInit() {
	c2BuffGeo = make([]float64, attributes.EndStatType)
	c2BuffGeo[attributes.CR] = 0.3

	c2BuffPyro = make([]float64, attributes.EndStatType)
	c2BuffPyro[attributes.ATKP] = 0.4

	c2BuffHydro = make([]float64, attributes.EndStatType)
	c2BuffHydro[attributes.HPP] = 0.4

	c2BuffCryo = make([]float64, attributes.EndStatType)
	c2BuffCryo[attributes.CD] = 0.5
}

func (c *char) c2buff() {
	if c.Base.Cons < 2 {
		return
	}
	for _, ch := range c.Core.Player.Chars() {
		switch ch.Base.Element {
		case attributes.Geo:
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.CR,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffGeo, true
					}
					return nil, false
				},
			})
		case attributes.Pyro:
			c.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffPyro, true
					}
					return nil, false
				},
			})
		}
	}
}

func (c *char) c2() {

}
