package linnea

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	c1Key     = "linnea-c1"
	c2Key     = "linnea-c2"
	c4Key     = "linnea-c4"
	c4KeySelf = "linnea-c4-self"
	c6Key     = "linnea-c6"
)

var lcrContributorMult = []float64{1.0, 1.0 / 2.0, 1.0 / 12.0, 1.0 / 12.0}

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		if !c.StatusIsActive(c1Key) {
			c.c1Stacks = 0
		}
		c.AddStatus(c1Key, 10*60, true)
		stacks := c.c6C1Stacks()
		c.c1Stacks = min(c.c1Stacks+stacks, 18)
		return false
	}, "linnea-c1")

	// this is emitted up to 4 times per attack, but we should only consume 1 stack per "attack"
	// so we only add damage in OnLunarReactionAttack, but we will consume in OnEnemyDamage
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagReactionLunarCrystallize:
		default:
			return false
		}
		if !c.StatusIsActive(c1Key) {
			return false
		}

		if c.c1Stacks == 0 {
			return false
		}

		maxStacks := 1
		scaling := 0.75

		c6stacks, c6scale := c.c6C1Mult()
		maxStacks *= c6stacks
		scaling *= c6scale

		def := c.TotalDef(false)
		stacks := min(c.c1Stacks, maxStacks)
		amt := def * scaling * float64(stacks)
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Linnea C1 proc dmg added to contribution", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("Field Catalog stacks left", c.c1Stacks)
		}
		atk.Info.FlatDmg += amt

		return false
	}, "linnea-c1-lcr-reaction")

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		maxStacks := 1
		scaling := 0.75
		if atk.Info.ActorIndex == c.Index() && atk.Info.Abil == skillMillionAbil {
			maxStacks = 5
			scaling = 1.5
		}

		if c.c1Stacks == 0 {
			return false
		}

		c6stacks, c6scale := c.c6C1Mult()
		maxStacks *= c6stacks
		scaling *= c6scale

		if !c.StatusIsActive(c1Key) {
			return false
		}

		stacks := min(c.c1Stacks, maxStacks)

		switch atk.Info.AttackTag {
		case attacks.AttackTagReactionLunarCrystallize:
			// we added the damage in OnLunarReactionAttack so we only need to reduce stack count here
			c.c1Stacks -= stacks
			return false
		case attacks.AttackTagDirectLunarCrystallize:
		default:
			return false
		}

		c.c1Stacks -= stacks
		def := c.TotalDef(false)
		amt := def * scaling * float64(stacks)
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Linnea C1 proc dmg add", glog.LogPreDamageMod, atk.Info.ActorIndex).
				Write("before", atk.Info.FlatDmg).
				Write("addition", amt).
				Write("Field Catalog stacks left", c.c1Stacks)
		}
		atk.Info.FlatDmg += amt

		return false
	}, "linnea-c1-dmg")
}

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	if !c.StatusIsActive(c1Key) {
		c.c1Stacks = 0
	}
	c.AddStatus(c1Key, 10*60, true)
	c.c1Stacks = min(c.c1Stacks+6, 18)
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.4
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		for _, char := range c.Core.Player.Chars() {
			switch char.Base.Element {
			case attributes.Geo:
			case attributes.Hydro:
			default:
				return false
			}
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c2Key, 8*60),
				AffectedStat: attributes.CD,
				Amount: func() ([]float64, bool) {
					return m, true
				},
			})
		}
		return false
	}, "linnea-c2")
}

func (c *char) c2MillionTonCDBonus(snap *info.Snapshot) {
	if c.Base.Cons < 2 {
		return
	}
	snap.Stats[attributes.CD] += 1.5
}

func (c *char) c2TriggerMoonDrift(ae *info.AttackEvent) {
	if c.Base.Cons < 2 {
		return
	}
	if c.Core.Player.GetMoonsignCount() < 2 {
		return
	}
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "Triggering C2 Moondrift Harmony")
	c.Core.Events.Emit(event.OnMoondriftHarmony, c.Core.Combat.PrimaryTarget(), ae)
	// c.Core.Log.NewEventBuildMsg(glog.LogElementEvent, c.Index(), "Linnea C2 Lunar Crystallize attack triggered")
	for _, delay := range []int{1, 4, 7} {
		c.Core.Tasks.Add(func() { c.doSingleLCrAttack() }, delay)
		if chance, ok := c.Core.Flags.Custom[reactable.LcrExtraHitOverride]; ok && c.Core.Rand.Float64() < chance {
			c.Core.Tasks.Add(func() { c.doSingleLCrAttack() }, delay)
		}
	}
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEFP] = 0.25
	c.Core.Events.Subscribe(event.OnMoondriftHarmony, func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(c4Key, 5*60),
				AffectedStat: attributes.DEFP,
				Amount: func() ([]float64, bool) {
					if c.Core.Player.Active() == char.Index() {
						return m, true
					}
					return nil, false
				},
			})
		}

		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(c4KeySelf, 5*60),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
		return false
	}, "linnea-c4")
}

func (c *char) c6C1Stacks() int {
	if c.Base.Cons < 6 {
		return 6
	}
	return 18
}

func (c *char) c6C1Mult() (int, float64) {
	if c.Base.Cons < 6 {
		return 1, 1.0
	}
	return 2, 1.5
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	if c.Core.Player.GetMoonsignCount() < 2 {
		return
	}

	amt := 0.25
	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) bool {
		atk := args[0].(*info.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagDirectLunarCrystallize {
			atk.Info.Elevation += amt
		}
		return false
	}, c6Key+"-direct")

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagReactionLunarCrystallize {
			atk.Info.Elevation += amt
		}
		return false
	}, c6Key+"-reaction")
}

func (c *char) doSingleLCrAttack() {
	contributions := []lcrContribution{}

	ap := combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key())

	// Do we need to make a new one for each character?
	ai := info.AttackInfo{
		DamageSrc:        c.Core.Combat.PrimaryTarget().Key(),
		Abil:             string(info.ReactionTypeLunarCrystallize),
		AttackTag:        attacks.AttackTagReactionLunarCrystallize,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Geo,
		IgnoreDefPercent: 1,
	}

	for charInd, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Geo:
		case attributes.Hydro:
		default:
			continue
		}

		ai.ActorIndex = charInd
		snap := char.Snapshot(&ai)

		ae := info.AttackEvent{
			Info:        ai,
			Pattern:     ap,
			SourceFrame: c.Core.F,
			Snapshot:    snap,
		}

		// Emit even so PreDamageMods can be applied to the individual LC contributions
		// Is there a way to collect these attackMods to show in logs?
		c.Core.Events.Emit(event.OnLunarReactionAttack, c.Core.Combat.PrimaryTarget(), &ae)

		em := ae.Snapshot.Stats[attributes.EM]
		cr := ae.Snapshot.Stats[attributes.CR]
		cd := ae.Snapshot.Stats[attributes.CD]

		flatdmg := 0.96 * combat.CalcLunarDmg(char.Base.Level, char, ae.Info, em)
		isCrit := false

		if c.Core.Rand.Float64() <= cr {
			flatdmg *= (1 + cd)
			isCrit = true
		}

		contributions = append(contributions, lcrContribution{flatdmg, isCrit, charInd, cr, cd, em})
	}

	if len(contributions) == 0 {
		return
	}

	slices.SortStableFunc(contributions, func(i, j lcrContribution) int {
		diff := j.dmg - i.dmg
		switch {
		case diff < 0:
			return -1
		case diff > 0:
			return 1
		default:
			return 0
		}
	})

	for i, contr := range contributions {
		c.Core.Combat.Log.NewEventBuildMsg(glog.LogElementEvent, contr.charInd, "lunarcrystallize contributor ", fmt.Sprint(i+1)).
			Write("target", c.Core.Combat.PrimaryTarget().Key()).
			Write("damage", &contr.dmg).
			Write("crit", &contr.isCrit).
			Write("mult", lcrContributorMult[i]).
			Write("cr", &contr.cr).
			Write("cd", &contr.cd).
			Write("em", &contr.em)

		ai.FlatDmg += contr.dmg * lcrContributorMult[i]
	}

	snap := info.Snapshot{}
	if contributions[0].isCrit {
		snap.Stats[attributes.CR] = 1.0
	}
	ai.ActorIndex = c.Index()
	c.Core.QueueAttackWithSnap(
		ai,
		snap,
		ap,
		0,
	)
}

type lcrContribution = struct {
	dmg     float64
	isCrit  bool
	charInd int
	cr      float64
	cd      float64
	em      float64
}
