package klee

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	a1IcdKey             = "a1-icd"
	a1SparkKey           = "a1-spark"
	boomboosterBurstKey  = "boombooster-burst"
	boomboosterNormalKey = "boombooster-normal"
	boomboosterSkillKey  = "boomboster-skill"
)

// When Jumpy Dumpty and Normal Attacks deal DMG, Klee has a 50% chance to obtain an Explosive Spark.
// This Explosive Spark is consumed by the next Charged Attack, which costs no Stamina and deals 50% increased DMG.
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
		// TODO: Does ICD change with magic?
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

func (c *char) getMagicCount() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		if c.IsMagic {
			count += 1
		}
	}
	return count
}

// Magic: Secret Rite
// Each time Klee deals DMG with her Elemental Skill, Elemental Burst, or Normal Attack, she gains a stack of Boom Booster.
// Each stack lasts for 20s and has its own independent timer. While Klee has 1/2/3 stacks, her Explosive Spark-enhanced
// Charged Attacks deal 115%/130%/150% of their original DMG.
func (c *char) magicInit() {
	if !c.IsMagic {
		return
	}
	if c.getMagicCount() < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		if c.Core.Player.Active() != c.Index() {
			return false
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
			c.AddStatus(boomboosterNormalKey, 60*20, true)
		case attacks.AttackTagElementalArt:
			c.AddStatus(boomboosterSkillKey, 60*20, true)
		case attacks.AttackTagElementalBurst:
			c.AddStatus(boomboosterBurstKey, 60*20, true)
		default:
			return false
		}

		return false
	}, "klee-boombooster")
}
