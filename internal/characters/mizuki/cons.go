package mizuki

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key               = "mizuki-c1"
	c1Interval          = 3.5 * 60
	c1Duration          = 3 * 60
	c1Multiplier        = 11.0
	c1Range             = 12
	c2Key               = "mizuki-c2"
	c2EMMultiplier      = 0.0004
	c2Interval          = 0.5 * 60
	c4EnergyGenerations = 4
	c4Key               = "mizuki-c4"
	c4Energy            = 5
	c6Key               = "mizuki-c6"
	c6CR                = 0.3
	c6CD                = 1.0
)

// When Yumemizuki Mizuki is in the Dreamdrifter state, she will continuously apply the "Twenty-Three Nights' Awaiting"
// effect to nearby opponents for 3s every 3.5s. When an opponent is affected by Anemo DMG-triggered Swirl reactions
// while the aforementioned effect is active, the effect will be canceled and this Swirl instance has its DMG against
// this opponent increased by 1100% of Mizuki's Elemental Mastery.
func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		e, ok := args[0].(*enemy.Enemy)
		atk := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		// Check if enemy has the debuff
		if !e.StatusIsActive(c1Key) {
			return
		}

		// Only on swirls. The swirl source does not matter, it can be either mizuki or another anemo char.
		switch atk.Info.AttackTag {
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlElectro:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlPyro:
		case attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		// do not proc on 0 DMG swirls (e.g. hydro AOE swirls or swirl ICD)
		if atk.Info.FlatDmg == 0 {
			return
		}

		additionalDmg := c1Multiplier * c.c1EM

		c.Core.Log.NewEvent("mizuki c1 proc", glog.LogPreDamageMod, atk.Info.ActorIndex).
			Write("before", atk.Info.FlatDmg).
			Write("addition", additionalDmg).
			Write("final", atk.Info.FlatDmg+additionalDmg)

		atk.Info.FlatDmg += additionalDmg
		atk.Info.Abil += " (Mizuki C1)"

		// Cancel the effect
		e.DeleteStatus(c1Key)

		if !c.revelation {
			return
		}

		// She also launches an additional attack on the same opponent,
		// dealing Anemo DMG equal to 1,100% of Yumemizuki Mizuki's Elemental Mastery.

		mult := 11.0
		if atk.Info.AttackTag == attacks.AttackTagReactionStellarSwirl {
			mult = 5.5
		}

		// this shouldn't swirl because it's right after an existing swirl?
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Mizuki C1",
			AttackTag:  attacks.AttackTagNone,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 25,
			Mult:       mult,
			UseEM:      true,
		}

		if c.isRadianceSSw() {
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Mult = 6
			ai.IgnoreDefPercent = 1
			ai.Durability = 0
		}

		ap := combat.NewCircleHitOnTarget(
			e,
			nil,
			5.5,
		)
		c.Core.QueueAttack(ai, ap, 3, 3)
	}, c1Key)
}

func (c *char) c1Task(src, hitmark int) {
	c.QueueCharTask(func() {
		if c.cloudSrc != src {
			return
		}
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return
		}

		c.c1EM = c.Stat(attributes.EM)
		area := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, c1Range)
		for _, target := range c.Core.Combat.EnemiesWithinArea(area, nil) {
			if e, ok := target.(*enemy.Enemy); ok {
				// is it even possible to verify if it is affected by hitlag?
				e.AddStatus(c1Key, c1Duration, true)
			}
		}
		c.c1Task(src, c1Interval)
	}, hitmark)
}

// When Yumemizuki Mizuki enters the Dreamdrifter state, every Elemental Mastery point she has will increase all nearby
// party members' Pyro, Hydro, Cryo, and Electro DMG Bonuses by 0.04% until the Dreamdrifter state ends.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2Buff = make([]float64, attributes.EndStatType)
	c.c2UpdateTask()

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		// TODO: Test whether this is indeed a static buff once we have C2
		char.AddStatMod(character.StatMod{
			Base: modifier.NewBase(c2Key, -1),
			Amount: func() []float64 {
				if !c.StatusIsActive(dreamDrifterStateKey) {
					return nil
				}
				return c.c2Buff
			},
		})
	}
}

func (c *char) c2UpdateTask() {
	if c.Base.Cons < 2 {
		return
	}

	c.QueueCharTask(func() {
		dmgBonus := c.NonExtraStat(attributes.EM) * c2EMMultiplier
		c.c2Buff[attributes.PyroP] = dmgBonus
		c.c2Buff[attributes.HydroP] = dmgBonus
		c.c2Buff[attributes.ElectroP] = dmgBonus
		c.c2Buff[attributes.CryoP] = dmgBonus

		c.c2UpdateTask()
	}, c2Interval)
}

func (c *char) c2OnSkill() {
	if c.Base.Cons < 2 {
		return
	}

	if !c.revelation {
		return
	}

	c.c2Src = c.Core.F
	c.c2Ticker(c.c2Src)
}

func (c *char) c2OnSkillExit() {
	if c.Base.Cons < 2 {
		return
	}

	if !c.revelation {
		return
	}

	c.c2Src = -1
}

func (c *char) c2Ticker(src int) {
	if !c.StatusIsActive(dreamDrifterStateKey) {
		return
	}

	if c.c2Src != src {
		return
	}

	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10)

	elems := []attributes.Element{attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo, attributes.Anemo}
	for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
		e, ok := e.(*enemy.Enemy)
		if !ok {
			continue
		}
		for _, elem := range elems {
			e.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag(c2Key+"-"+elem.String(), 1*60),
				Ele:   elem,
				Value: -0.20,
			})
		}

	}

	c.Core.Tasks.Add(func() { c.c2Ticker(src) }, 0.3*60)
}

// Picking up a Yumemi Style Special Snack from the Elemental Burst Anraku Secret Spring Therapy will both deal DMG
// and heal, and will restore 5 Energy to Yumemizuki Mizuki. Energy can be restored this way 4 times per Anraku
// Secret Spring Therapy duration.
func (c *char) c4() {
	if c.Base.Cons < 4 {
		return
	}

	if c.c4EnergyGenerationsRemaining > 0 {
		c.c4EnergyGenerationsRemaining--
		c.AddEnergy(c4Key, c4Energy)
	}

	if !c.revelation {
		return
	}

	lowestHPInd := c.Index()
	lowestHP := c.CurrentHPRatio()
	for i, char := range c.Core.Player.Chars() {
		if char.CurrentHPRatio() < lowestHP {
			lowestHPInd = i
			lowestHP = char.CurrentHPRatio()
		}
	}

	healAmt := 2.66 * c.Stat(attributes.EM)
	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  lowestHPInd,
		Message: snackHealName + " (C4)",
		Src:     healAmt,
		Bonus:   c.Stat(attributes.Heal),
	})
}

// While Yumemizuki Mizuki is in the Dreamdrifter state, Swirl DMG dealt by nearby party members can Crit,
// with CRIT Rate fixed at 30%, and CRIT DMG fixed at 100%.
func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		// Only on swirls. The swirl source does not matter, it can be either mizuki or other anemo char.
		switch ae.Info.AttackTag {
		case attacks.AttackTagSwirlPyro:
		case attacks.AttackTagSwirlCryo:
		case attacks.AttackTagSwirlHydro:
		case attacks.AttackTagSwirlElectro:
		default:
			return
		}

		// The effect is only when mizuki is in dreamDrifter state
		if !c.StatusIsActive(dreamDrifterStateKey) {
			return
		}

		// Crit rate/DMG is fixed to 30% CR and 100% CD
		ae.Snapshot.Stats[attributes.CR] = c6CR
		ae.Snapshot.Stats[attributes.CD] = c6CD
	}, c6Key)

	if !c.revelation {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c6Key, -1),
		Extra:        true,
		AffectedStat: attributes.NoStat,
		Amount: func() []float64 {
			scaling := max(c.NonExtraStat(attributes.EM)-500, 0)
			m[attributes.CR] = min(scaling*0.0004, 0.2)
			m[attributes.CD] = min(scaling*0.0016, 0.8)
			return m
		},
	})

	m2 := make([]float64, attributes.EndStatType)
	m2[attributes.CR] = 0.10
	m2[attributes.CD] = 0.20
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c6Key+"ssw-cr", -1),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				if atk.Info.AttackTag != attacks.AttackTagDirectStellarSwirl {
					return nil
				}

				if !c.StatusIsActive(dreamDrifterStateKey) {
					return nil
				}
				return m2
			},
		})
	}

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if ae.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		ae.Snapshot.Stats[attributes.CR] += 0.1
		ae.Snapshot.Stats[attributes.CD] += 0.2
	}, c6Key+"-lunar-reaction")
}
