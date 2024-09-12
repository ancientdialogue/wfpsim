package xilonen

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const skillStart = 13
const skillHitmarks = 17
const particleICDKey = "xilonen-particle-icd"
const samplerShredKey = "xilonen-e-shred"
const activeSamplerKey = "xilonen-samplers-activated"
const maxNightsoulPoints = 90

func init() {
	skillFrames = frames.InitAbilSlice(52)
	skillFrames[action.ActionAttack] = 29
	skillFrames[action.ActionBurst] = 28
	skillFrames[action.ActionDash] = 37
	skillFrames[action.ActionJump] = 37
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	if c.nightsoulState.Points() <= 0 {
		c.exitNightsoul()
	}
}

func (c *char) enterNightsoul() {
	c.nightsoulState.EnterBlessing(45)
	c.nightsoulSrc = c.Core.F
	interval := int(12.0 * c.c1IntervalMod())
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.nightsoulSrc), interval)
	c.NormalHitNum = rollerHitNum

	src := c.nightsoulSrc
	duration := int(9 * 60 * c.c1DurMod())
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		c.exitNightsoul()
	}, duration)

	// Don't queue the task if C2 or higher
	if c.Base.Cons < 2 {
		c.activeGeoSampler(src)
	}
}

func (c *char) exitNightsoul() {
	c.nightsoulState.ExitBlessing()
	c.nightsoulSrc = -1
	c.NormalHitNum = normalHitNum
	c.SetCDWithDelay(action.ActionSkill, 7*60, 0)
	c.NormalCounter = 0
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}

		if c.nightsoulState.Points() <= 0.001 {
			return
		}

		// RODO: is this check needed? The nightsoulSrc gets reset on on exiting NS state
		// if !c.nightsoulState.HasBlessing() {
		// 	return
		// }

		c.reduceNightsoulPoints(1)
		interval := int(12.0 * c.c1IntervalMod())
		// reduce 1 point per 12f, which is 5 per second
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), interval)
	}
}

func (c *char) applySamplerShred(ele attributes.Element, enemies []combat.Enemy) {
	for _, e := range enemies {
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag(fmt.Sprintf("%v-%v", samplerShredKey, ele.String()), 60),
			Ele:   ele,
			Value: -skillShred[c.TalentLvlSkill()],
		})
	}
}

func (c *char) activeGeoSampler(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
		c.applySamplerShred(attributes.Geo, enemies)

		// TODO: how often does this apply?
		c.QueueCharTask(c.activeGeoSampler(src), 18)
	}
}

func (c *char) activeSamplers(src int) func() {
	return func() {
		if c.sampleSrc != src {
			return
		}

		if !c.StatusIsActive(activeSamplerKey) {
			return
		}

		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)

		for _, ele := range c.shredElements {
			c.applySamplerShred(ele, enemies)
		}

		// TODO: how often does this apply?
		c.QueueCharTask(c.activeSamplers(src), 18)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.exitNightsoul()
		return action.Info{
			Frames:          func(_ action.Action) int { return 1 },
			AnimationLength: 1,
			CanQueueAfter:   1, // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Yohual's Scratch",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Geo,
		Durability:         25,
		HitlagHaltFrames:   0.02 * 60,
		HitlagFactor:       0.01,
		Mult:               skillDMG[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1},
		3,
	)
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
	}, skillHitmarks)

	c.QueueCharTask(func() {
		c.enterNightsoul()
	}, skillStart)

	c.c4()
	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Geo, c.ParticleDelay)
}