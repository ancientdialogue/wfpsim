package venti

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstKey   = "stormeye"
	burstStart = 94
)

func init() {
	burstFrames = frames.InitAbilSlice(95) // Q -> N1/CA/E/D
	burstFrames[action.ActionJump] = 94    // Q -> J
	burstFrames[action.ActionSwap] = 93    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// reset location
	c.qAbsorb = attributes.NoElement
	player := c.Core.Combat.Player()
	c.qPos = info.CalcOffsetPoint(player.Pos(), info.Point{Y: 5}, player.Direction())
	c.absorbCheckLocation = combat.NewBoxHitOnTarget(c.qPos, info.Point{Y: -1}, 2.5, 2.5)

	// 8 second duration, tick every .4 second
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Wind's Grand Ode",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurstAnemo,
		ICDGroup:   attacks.ICDGroupVenti,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burstDot[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.qPos, nil, 4)

	c.aiAbsorb = ai
	c.aiAbsorb.Abil = "Wind's Grand Ode (Absorbed)"
	c.aiAbsorb.Mult = burstAbsorbDot[c.TalentLvlBurst()]
	c.aiAbsorb.Element = attributes.NoElement

	// snapshot is around cd frame and 1st tick?
	var snap info.Snapshot
	c.Core.Tasks.Add(func() {
		snap = c.Snapshot(&ai)
		c.snapAbsorb = c.Snapshot(&c.aiAbsorb)

		c.qAbsorbBonusTicks = 0

		c.c2OnBurst()
		c.c4OnSkillBurst()
	}, 104)

	cb := c.c6(attributes.Anemo)

	// starts at 106 with 24f interval between ticks. 20 total
	c.qSrc = c.Core.F
	c.Core.Tasks.Add(c.burstTicks(c.Core.F, ai, &snap, ap, cb), 106)
	c.Core.Tasks.Add(func() { c.AddStatus(burstKey, 8*60, false) }, 93)

	// Infusion usually occurs after 4 ticks of anemo according to KQM library
	c.Core.Tasks.Add(c.absorbCheckQ(c.Core.F, 0, int((480-24*4)/18)), 106+24*3)

	// Ensure that the hex extension cap gets reset
	c.hexereiBurstExtCount = 0

	if c.Base.Ascension >= 4 {
		c.Core.Tasks.Add(c.a4, 480+burstStart)
	}

	c.SetCDWithDelay(action.ActionBurst, 15*60, 81)
	c.ConsumeEnergy(84)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTicks(src int, ai info.AttackInfo, snap *info.Snapshot, ap info.AttackPattern, cb info.AttackCBFunc) func() {
	return func() {
		if c.qSrc != src {
			return
		}
		if !c.StatusIsActive(burstKey) {
			return
		}
		c.Core.QueueAttackWithSnap(ai, *snap, ap, 0, cb)
		ai.Mult = burstDot[c.TalentLvlBurst()] * c.hexereiBurstBuff()
		c.Core.Tasks.Add(c.burstTicks(src, ai, snap, ap, cb), 24)
	}
}

// TODO: how does this work with the bug and the extension? just add more ticks?
func (c *char) burstAbsorbedTicks(count int, ai info.AttackInfo, snap info.Snapshot, ap info.AttackPattern, cb info.AttackCBFunc) func() {
	// ticks at 24f. 15 total
	return func() {
		ai.Mult = burstDot[c.TalentLvlBurst()] * c.hexereiBurstBuff()
		c.Core.QueueAttackWithSnap(c.aiAbsorb, c.snapAbsorb, ap, 0, cb)
		if count+c.qAbsorbBonusTicks <= 0 {
			c.aiAbsorb.Element = attributes.NoElement
			return
		}

		c.Core.Tasks.Add(c.burstAbsorbedTicks(count-1, ai, snap, ap, cb), 24)
	}
}

func (c *char) absorbCheckQ(src, count, maxcount int) func() {
	return func() {
		if count == maxcount {
			return
		}
		c.qAbsorb = c.Core.Combat.AbsorbCheck(c.Index(), c.absorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if c.qAbsorb != attributes.NoElement {
			c.aiAbsorb.Element = c.qAbsorb
			switch c.qAbsorb {
			case attributes.Pyro:
				c.aiAbsorb.ICDTag = attacks.ICDTagElementalBurstPyro
			case attributes.Hydro:
				c.aiAbsorb.ICDTag = attacks.ICDTagElementalBurstHydro
			case attributes.Electro:
				c.aiAbsorb.ICDTag = attacks.ICDTagElementalBurstElectro
			case attributes.Cryo:
				c.aiAbsorb.ICDTag = attacks.ICDTagElementalBurstCryo
			}
			// trigger dmg ticks here

			cb := c.c6(c.qAbsorb)

			ap := combat.NewCircleHitOnTarget(c.qPos, nil, 6)
			c.Core.Tasks.Add(c.burstAbsorbedTicks(15, c.aiAbsorb, c.snapAbsorb, ap, cb), 0)
			return
		}
		// otherwise queue up
		c.Core.Tasks.Add(c.absorbCheckQ(src, count+1, maxcount), 18)
	}
}
