package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	a1IcdKey           = "a1-icd"
	a1SparkKey         = "a1-spark"
	boomBadgeBurstKey  = "boombadge-burst"
	boomBadgeNormalKey = "boombadge-normal"
	boomBadgeSkillKey  = "boombadge-skill"
)

// When Jumpy Dumpty and Normal Attacks deal DMG, Klee has a 50% chance to obtain an Explosive Spark.
// This Explosive Spark is consumed by the next Charged Attack, which costs no Stamina and deals 50% increased DMG.

// Hexerei bonus:
// Using her Elemental Skill Jumpy Dumpty or Elemental Burst Sparks 'n' Splash also grants 1 Explosive Spark.
// Klee can hold up to 3 Explosive Sparks. While holding an Explosive Spark, Klee's Charged Attack will consume
// 1 Spark to perform a special stamina-free Charged Attack, Boom-Boom Strike, dealing 50% increased DMG.
func (c *char) makeA1CB() info.AttackCBFunc {
	if c.Base.Ascension < 1 {
		return nil
	}
	return func(a info.AttackCB) {
		if c.StatusIsActive(a1IcdKey) {
			return
		}
		if c.Core.Rand.Float64() < 0.5 {
			return
		}
		// TODO: Does ICD change with hexerei?
		c.AddStatus(a1IcdKey, 60*5, true)
		c.addSpark()
	}
}

func (c *char) addSpark() {
	if c.Base.Ascension < 1 {
		return
	}

	previous := c.a1CurrentStack
	c.a1CurrentStack++
	// TODO: Does status refresh at max stacks?
	if c.a1CurrentStack > c.a1MaxStack {
		c.a1CurrentStack = c.a1MaxStack
		return
	}
	c.AddStatus(a1SparkKey, 60*30, true)
	c.Core.Log.NewEvent("adding spark stack", glog.LogCharacterEvent, c.Index()).
		Write("previous", previous).
		Write("new", c.a1CurrentStack)
}

const a4ICDKey = "klee-a4-icd"

// When Klee's Charged Attack results in a CRIT Hit, all party members gain 2 Elemental Energy.
func (c *char) makeA4CB() info.AttackCBFunc {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		if !a.IsCrit {
			return
		}
		if c.StatusIsActive(a4ICDKey) {
			return
		}
		c.AddStatus(a4ICDKey, 0.6*60, true)
		for _, x := range c.Core.Player.Chars() {
			x.AddEnergy("klee-a4", 2)
		}
	}
}

// Hexerei: Secret Rite
// When Klee deals DMG with Normal Attacks, Elemental Skill, or Elemental Burst, she gains 1 Boom Badge, lasting 20s.
// Each type of attack can grant at most 1 Boom Badge this way, and each badge has its own independent timer.
// While Klee has 1/2/3 Boom Badges, her special Charged Attack Boom-Boom Strike deals 115%/130%/150% of its original DMG.
func (c *char) hexereiInit() {
	if !c.IsHexerei {
		return
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		if c.Core.Player.Active() != c.Index() {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			c.AddStatus(boomBadgeNormalKey, 60*20, true)
		case attacks.AttackTagElementalArt:
			c.AddStatus(boomBadgeSkillKey, 60*20, true)
		case attacks.AttackTagElementalBurst:
			c.AddStatus(boomBadgeBurstKey, 60*20, true)
		default:
			return
		}
	}, "klee-boombadge")
}
