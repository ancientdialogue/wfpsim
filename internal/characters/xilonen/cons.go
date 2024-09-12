package xilonen

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2key = "xilonen-c2"
const c4key = "xilonen-c4"

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
	c2BuffGeo[attributes.DmgP] = 0.4

	c2BuffPyro = make([]float64, attributes.EndStatType)
	c2BuffPyro[attributes.ATKP] = 0.4

	c2BuffHydro = make([]float64, attributes.EndStatType)
	c2BuffHydro[attributes.HPP] = 0.4

	c2BuffCryo = make([]float64, attributes.EndStatType)
	c2BuffCryo[attributes.CD] = 0.5
}

func (c *char) c2buff() {
	for _, ch := range c.Core.Player.Chars() {
		switch ch.Base.Element {
		case attributes.Geo:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.DmgP,
				Amount: func() ([]float64, bool) {
					// geo is always active
					return c2BuffGeo, true
				},
			})
		case attributes.Pyro:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.ATKP,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffPyro, true
					}
					return nil, false
				},
			})
		case attributes.Hydro:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.HPP,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffHydro, true
					}
					return nil, false
				},
			})
		case attributes.Cryo:
			ch.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2key, -1),
				AffectedStat: attributes.CD,
				Amount: func() ([]float64, bool) {
					if c.StatusIsActive(activeSamplerKey) {
						return c2BuffCryo, true
					}
					return nil, false
				},
			})
		}
	}
}

func (c *char) c2GeoSampler() func() {
	return func() {
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
		c.applySamplerShred(attributes.Geo, enemies)

		// TODO: how often does this apply?
		c.QueueCharTask(c.c2GeoSampler(), 30)
	}
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2GeoSampler()()
	c.c2buff()
}

func (c *char) c2electro() {
	if c.Base.Cons < 2 {
		return
	}
	for _, ch := range c.Core.Player.Chars() {
		if ch.Base.Element == attributes.Electro {
			ch.AddEnergy(c2key, 20)
			ch.ReduceActionCooldown(action.ActionBurst, 5*60)
		}
	}
}

func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		char.AddStatus(c4key, 15*60, true) // 15 sec duration
		char.SetTag(c4key, 6)              // 6 c4 stacks
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		char := c.Core.Player.ByIndex(atk.Info.ActorIndex)

		if !char.StatusIsActive(c4key) {
			return false
		}

		if char.Tags[c4key] > 0 {
			amt := 0.65 * c.TotalDef()
			char.Tags[c4key]--

			c.Core.Log.NewEvent("Xilonen c4 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("effect_ends_at", c.StatusExpiry(c4key)).
				Write("c4_left", c.Tags[c4key])

			atk.Info.FlatDmg += amt
		}

		return false
	}, fmt.Sprintf("%s-hook", c4key))
}
