package chasca

import (
	"math/rand"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	skillStateKey       = "skill-state-icd"
	particleICD         = 9999 * 60
	particleICDKey      = "chasca-particle-icd"
	maxShadowhuntShells = 6
)

func init() {
	skillFrames = frames.InitAbilSlice(30)
}
func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0.00001 {
		c.cancelNightsoul()
	}
}

func (c *char) cancelNightsoul() {
	c.nightsoulState.ExitBlessing()
	c.SetCDWithDelay(action.ActionSkill, 6.5*60, 0)
	c.ResetActionCooldown(action.ActionAim)
	c.nightsoulSrc = -1
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}

		if !c.nightsoulState.HasBlessing() {
			return
		}

		c.reduceNightsoulPoints(8)

		// reduce 8 point per 60f
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), 60)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.cancelNightsoul()
		return action.Info{
			Frames:          func(_ action.Action) int { return 1 },
			AnimationLength: 1,
			CanQueueAfter:   1, // earliest cancel
			State:           action.SkillState,
		}, nil
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Resonance DMG",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			nil,
			2.5,
		),
		0,
		0,
	)
	c.QueueCharTask(func() {
		c.nightsoulState.EnterBlessing(80)
		c.DeleteStatus(particleICDKey)
		c.nightsoulSrc = c.Core.F
		c.QueueCharTask(c.nightsoulPointReduceFunc(c.nightsoulSrc), 8)
	}, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionCharge], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) loadShadowhuntShells(shellNum int) {
	shellNum = max(c.maxshadowhuntShell, shellNum) // max 6 shell
	for i := range c.shadowhuntShells {
		c.shadowhuntShells[i] = attributes.Anemo
	}
	rand.Shuffle(len(c.conversionElements), func(i, j int) {
		c.conversionElements[i], c.conversionElements[j] = c.conversionElements[j], c.conversionElements[i]
	})

	for i := 0; i < shellNum; i++ {
		if i < len(c.conversionElements) {
			c.shadowhuntShells[shellNum-1-i] = c.conversionElements[i]
		} else {
			c.shadowhuntShells[shellNum-1-i] = attributes.Anemo
		}
	}
	c.a1()
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	particleCount := 5.0
	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Hydro, c.ParticleDelay)
}
