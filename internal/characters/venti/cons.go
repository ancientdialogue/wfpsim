package venti

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c2BuffKey = "c2-skill-dmg"

// C1:
// Fires 2 additional arrows per Aimed Shot, each dealing 33% of the original arrow's DMG.
func (c *char) c1Charge(ai info.AttackInfo, hitmark, travel int) {
	if c.Base.Cons < 1 {
		return
	}
	ai.Abil += " (C1)"
	ai.Mult /= 3.0
	for range 2 {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				info.Point{Y: -0.5},
				0.1,
				1,
			),
			hitmark,
			hitmark+travel,
		)
	}
}

// C1:
func (c *char) c1Normal(ai info.AttackInfo, hitmark, travel int) {
	if c.Base.Cons < 1 {
		return
	}
	ai.Abil += " (C1)"
	ai.Mult *= 0.2
	for range 2 {
		c.Core.QueueAttack(
			ai,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				info.Point{Y: -0.5},
				0.1,
				1,
			),
			hitmark,
			hitmark+travel,
		)
	}
}

// C2:
// Skyward Sonnet decreases opponents' Anemo RES and Physical RES by 12% for 10s.
// Opponents launched by Skyward Sonnet suffer an additional 12% Anemo RES and Physical RES decrease while airborne.
// TODO: the airborne part isn't implemented
func (c *char) c2SkillCB(a info.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	e, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-anemo", 600),
		Ele:   attributes.Anemo,
		Value: -0.24,
	})
	e.AddResistMod(info.ResistMod{
		Base:  modifier.NewBaseWithHitlag("venti-c2-phys", 600),
		Ele:   attributes.Physical,
		Value: -0.24,
	})
}

func (c *char) c2OnSkillTap() float64 {
	if c.Base.Cons < 2 {
		return 1.0
	}

	if !c.StatusIsActive(c2BuffKey) {
		return 1.0
	}
	c.DeleteStatus(c2BuffKey)
	return 3.0
}

func (c *char) c2OnBurst() {
	if c.Base.Cons < 2 {
		return
	}

	c.AddStatus(c2BuffKey, 15*60, true)
	c.ResetActionCooldown(action.ActionSkill)
}

func (c *char) c2NaCB(a info.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	_, ok := a.Target.(*enemy.Enemy)
	if !ok {
		return
	}

	if c.Core.Rand.Float64() > 0.25 {
		return
	}

	c.AddStatus(c2BuffKey, 15*60, true)
}

// C4:
// When Venti picks up an Elemental Orb or Particle, he receives a 25% Anemo DMG Bonus for 10s.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.c4bonus = make([]float64, attributes.EndStatType)
	c.c4bonus[attributes.AnemoP] = 0.25
}

func (c *char) c4OnSkillBurst() {
	if c.Base.Cons < 4 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		if c.Index() == char.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
			AffectedStat: attributes.AnemoP,
			Amount: func() []float64 {
				if char.Index() != c.Core.Player.Active() {
					return nil
				}
				return c.c4bonus
			},
		})
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
		AffectedStat: attributes.AnemoP,
		Amount: func() []float64 {
			return c.c4bonus
		},
	})

	// c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...any) bool {
	// 	// only trigger if Venti catches the particle
	// 	if c.Core.Player.Active() != c.Index() {
	// 		return false
	// 	}
	// 	// apply C4 to Venti
	// 	c.AddStatMod(character.StatMod{
	// 		Base:         modifier.NewBaseWithHitlag("venti-c4", 600),
	// 		AffectedStat: attributes.AnemoP,
	// 		Amount: func() ([]float64, bool) {
	// 			return c.c4bonus, true
	// 		},
	// 	})
	// 	return false
	// }, "venti-c4")
}

// C6:
// Targets who take DMG from Wind's Grand Ode have their Anemo RES decreased by 20%.
// If an Elemental Absorption occurred, then their RES towards the corresponding Element is also decreased by 20%.
func (c *char) c6(ele attributes.Element) func(a info.AttackCB) {
	if c.Base.Cons < 6 {
		return nil
	}

	return func(a info.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("venti-c6-"+ele.String(), 600),
			Ele:   ele,
			Value: -0.20,
		})
	}
}

// C6:
// Targets who take DMG from Wind's Grand Ode have their Anemo RES decreased by 20%.
// If an Elemental Absorption occurred, then their RES towards the corresponding Element is also decreased by 20%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c6bonus := make([]float64, attributes.EndStatType)
	c6bonus[attributes.CD] = 1.0

	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("venti-c6-cdmg", -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil
			}
			if !e.ResistModIsActive("venti-c6-" + attributes.Anemo.String()) {
				return nil
			}
			return c6bonus
		},
	})
}
