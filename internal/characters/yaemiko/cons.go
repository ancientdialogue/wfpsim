package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "yae-c1"
	c2Key = "yae-c2"
)

var c2BuffVal = []float64{0, 60, 90, 120, 200}

// After Yae Miko triggers a Superconduct or Stellar-Conduct reaction, or deals Stellar-Conduct DMG,
// nearby party members will gain a 50% Electro DMG and Stellar-Conduct DMG Bonus for 10s. Triggering
// these reactions again will refresh the duration of the DMG bonuses.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	if !c.revelation {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ElectroP] = 0.5

	gainBuffs := func() {
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c1Key+"-electro", 10*60),
				AffectedStat: attributes.ElectroP,
				Amount: func() []float64 {
					return m
				},
			})
			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBase(c1Key+"-ssc", 10*60),
				Amount: func(ai info.AttackInfo) float64 {
					if ai.AttackTag == attacks.AttackTagDirectStellarConduct {
						return 0.5
					}
					return 0
				},
			})
		}
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		gainBuffs()
	}

	hookDmg := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagDirectStellarConduct {
			return
		}

		gainBuffs()
	}

	c.Core.Events.Subscribe(event.OnSuperconduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnEnemyDamage, hookDmg, c1Key)
}

// Additionally, when there are Sesshou Sakura in the field, Yae Miko's and your current active
// character's Elemental Mastery will also be increased by 60/90/120/200 points, depending on
// the level of the Sesshou Sakura.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	if !c.revelation {
		return
	}
	m := make([]float64, attributes.EndStatType)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c2Key, 10*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				level := c.sakuraLevelCheck()
				if level == 0 {
					return nil
				}
				m[attributes.EM] = c2BuffVal[level]
				return m
			},
		})
	}
}

// When Sesshou Sakura lightning hits opponents, the Electro DMG Bonus of all nearby party members is increased by 20% for 5s.
func (c *char) c4() {
	// TODO: does this trigger for yaemiko too? assuming it does
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("yaemiko-c4", 5*60),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				return c.c4buff
			},
		})
	}
}
