package ororon

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const a1NSBurstKey = "ororon-a1-ns-burst"
const a1ElectroHydroKey = "ororon-a1-electro-hydro"
const a1ECKey = "ororon-a1-ec"
const a1NSKey = "ororon-a1-ns"

const a1OnSkillKey = "ororon-a1"
const a1GainIcdKey = "ororon-a1-gain-icd"
const a1DamageIcdKey = "ororon-a1-dmg-icd"
const a1Abil = "Hypersense"

const a4Key = "ororon-a4"
const a4IcdKey = "ororon-a4-icd"

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		c.nightsoulState.GeneratePoints(40)
		return false
	}, a1NSBurstKey)

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		// ignores ororon himself
		if atk.Info.ActorIndex == c.Index {
			return false
		}

		switch c.Core.Player.Chars()[atk.Info.ActorIndex].Base.Element {
		case attributes.Hydro:
		case attributes.Electro:
		default:
			return false
		}

		if !c.StatusIsActive(a1OnSkillKey) {
			return false
		}
		if c.StatusIsActive(a1GainIcdKey) {
			return false
		}
		c.AddStatus(a1GainIcdKey, 0.3*60, true)

		c.nightsoulState.GeneratePoints(5)
		c.Tags[a1ElectroHydroKey]++
		if c.Tags[a1ElectroHydroKey] >= 10 {
			c.DeleteStatus(a1OnSkillKey)
		}
		return false
	}, a1ElectroHydroKey)

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		if atk.Info.Abil != "electrocharged" {
			return false
		}

		c.a1NightSoulAttack()

		return false
	}, a1ECKey)

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)

		// ignores ororon himself
		if atk.Info.ActorIndex == c.Index {
			return false
		}

		if !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagNightsoul) {
			return false
		}

		c.a1NightSoulAttack()

		return false
	}, a1NSKey)
}

func (c *char) a1NightSoulAttack() {
	if c.nightsoulState.Points() < 10 {
		return
	}
	if c.StatusIsActive(a1DamageIcdKey) {
		return
	}
	c.AddStatus(a1DamageIcdKey, 1.8*60, true)
	if !c.nightsoulState.HasBlessing() {
		c.a1EnterBlessing()
	}
	c.nightsoulState.ConsumePoints(10)
	c.hypersense(1.6)
}

func (c *char) hypersense(mult float64) {
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           a1Abil,
		AttackTag:      attacks.AttackTagNone,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           mult,
	}

	enemies := []targets.TargetKey{c.Core.Combat.PrimaryTarget().Key()}
	for i := 0; len(enemies) < 4 && i < c.Core.Combat.EnemyCount(); i++ {
		newKey := c.Core.Combat.Enemies()[i].Key()
		if newKey == c.Core.Combat.PrimaryTarget().Key() {
			continue
		}
		enemies = append(enemies, newKey)
	}
	snap := c.Snapshot(&ai)
	for _, e := range enemies {
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewSingleTargetHit(e),
			3,
			c.makeC6cb(),
		)
	}
}

func (c *char) a1EnterBlessing() {
	c.nightsoulState.EnterBlessing(c.nightsoulState.Points())
	c.QueueCharTask(c.nightsoulState.ExitBlessing, 6*60)
}

func (c *char) a1OnSkill() {
	if c.Base.Ascension < 1 {
		return
	}
	c.AddStatus(a1OnSkillKey, 15*60, true)
	c.SetTag(a1OnSkillKey, 0)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}

		switch atk.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		default:
			return false
		}

		if !c.StatusIsActive(a4Key) {
			return false
		}

		if c.StatusIsActive(a4IcdKey) {
			return false
		}

		c.AddStatus(a4IcdKey, 60, true)

		c.Core.Player.ActiveChar().AddEnergy(a4Key, 3)

		if c.Core.Player.Active() != c.Index {
			c.AddEnergy(a4Key, 3)
		}

		c.Tags[a4Key]++
		if c.Tags[a4Key] >= 3 {
			c.DeleteStatus(a4Key)
		}
		return false
	}, a4Key)
}

func (c *char) makeA4cb() func(combat.AttackCB) {
	if c.Base.Ascension < 4 {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		c.AddStatus(a4Key, 15*60, true)
		c.SetTag(a4Key, 0)
	}
}
