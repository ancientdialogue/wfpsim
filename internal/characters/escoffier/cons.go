package escoffier

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Dur = 15 * 60

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Base.Ascension < 4 {
		return
	}
	c.c1Active = false
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Cryo, attributes.Hydro:
		default:
			return
		}
	}
	c.c1Active = true
	c.c1Buff = make([]float64, attributes.EndStatType)
	c.c1Buff[attributes.CD] = 0.60
}

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Base.Ascension < 4 {
		return
	}
	if !c.c1Active {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		// TODO: check if this buff is affected by hitlag on characters or hitlag on escoffier
		// Currently assuming affected by hitlag on characters
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag("escoffier-c1", c1Dur),
			Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if ae.Info.Element != attributes.Cryo {
					return nil, false
				}
				return c.c1Buff, true
			},
		})
	}
}
