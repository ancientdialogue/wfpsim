package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1Key = "klee-c1-atk%"

func (c *char) c1(delay int) {
	if c.Base.Cons < 1 {
		return
	}
	// 0.1 base change, + 0.08 every failure
	if c.Core.Rand.Float64() > c.c1Chance {
		// failed
		c.c1Chance += 0.08
		return
	}
	c.c1Chance = 0.1

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Sparks'n'Splash (C1)",
		AttackTag:          attacks.AttackTagElementalBurst,
		ICDTag:             attacks.ICDTagElementalBurst,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         25,
		Mult:               1.2 * burst[c.TalentLvlBurst()],
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	// TODO: should center on target hit by attack that triggered c1
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5), 0, delay)
	c.Core.Log.NewEvent("c1 triggered", glog.LogCharacterEvent, c.Index())

	// Hexerei bonus:
	// Additionally, bombarding opponents with sparks increases Klee's ATK by 60% for 12s.
	if !c.IsHexerei {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.6
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c1Key, 12*60),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	})
}

// Being hit by Jumpy Dumpty's mines decreases opponents' DEF by 23% for 10s.
func (c *char) makeC2CB(mine bool) info.AttackCBFunc {
	return func(a info.AttackCB) {
		if c.Base.Cons < 2 {
			return
		}
		// Hexerei bonus:
		// Dealing DMG to opponents with Klee's Elemental Skill decreases their DEF by 23% for 10s.
		if !mine && !c.IsHexerei {
			return
		}
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		e.AddDefMod(info.DefMod{
			Base:  modifier.NewBaseWithHitlag("kleec2", 10*60),
			Value: -0.233,
		})
	}
}

func (c *char) triggerC4() {
	if c.Base.Cons < 4 {
		return
	}
	activeMult := 1.0
	if c.IsHexerei && c.Core.Player.Active() == c.Index() {
		activeMult = 2.0
	}
	// blow up
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Sparkly Explosion (C4)",
		AttackTag:          attacks.AttackTagNone,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeDefault,
		Element:            attributes.Pyro,
		Durability:         50,
		Mult:               5.55 * activeMult,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	// TODO: delay?
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0)
}
