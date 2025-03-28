package escoffier

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Dur = 15 * 60
const c1Key = "escoffier-c1"
const c2Key = "escoffier-c2"
const c2Per = 2.4
const c2Dur = 15 * 60

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
			Base: modifier.NewBaseWithHitlag(c1Key, c1Dur),
			Amount: func(ae *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
				if ae.Info.Element != attributes.Cryo {
					return nil, false
				}
				return c.c1Buff, true
			},
		})
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if c.Index == atk.Info.ActorIndex {
			return false
		}
		if c.Core.Player.Active() != atk.Info.ActorIndex {
			return false
		}
		if !c.StatusIsActive(c2Key) {
			return false
		}
		if c.c2Count > 0 {
			amt := c.TotalAtk() * c2Per
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Escoffier C2 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
					Write("before", atk.Info.FlatDmg).
					Write("addition", amt).
					Write("Cold Dishes left", c.c2Count)
			}
			atk.Info.FlatDmg += amt
			c.c2Count--
		}
		return false
	}, c2Key+"-on-dmg")
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.AddStatus(c2Key, c2Dur, true)
	c.c2Count = 5
}
