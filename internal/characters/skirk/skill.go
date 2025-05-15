package skirk

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const (
	maxSerpentsSubtlety = 100
	skillHitmark        = 28
	skillKey            = "seven-phase-slash"
	particleICDKey      = "skirk-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(48) // E -> N1
	skillFrames[action.ActionBurst] = 48   // E -> Q
	skillFrames[action.ActionDash] = 25    // E -> D
	skillFrames[action.ActionJump] = 26    // E -> J
	skillFrames[action.ActionSwap] = 49    // E -> Swap
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	h := p["hold"]
	if h > 0 {
		return c.skillHold(p)
	}
	return c.skillTap()
}

func (c *char) skillTap() (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		c.exitSkillState(c.skillSrc)
	} else {
		c.QueueCharTask(func() { c.enterSkillState() }, 26)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel is before skillHitmark
		State:           action.SkillState,
	}, nil
}

func (c *char) enterSkillState() {
	c.skillSrc = c.Core.F
	c.AddStatus(skillKey, 12.5*60, false)
	c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill", 45.0)
	c.c2OnSkill()
	c.serpentsReduceTask(c.skillSrc)
	src := c.skillSrc
	c.Core.Tasks.Add(func() { c.exitSkillState(src) }, 12.5*60)
}

func (c *char) exitSkillState(src int) {
	if c.skillSrc != src {
		return
	}
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index, "exit skirk skill").Write("src", src)
	c.skillSrc = -1
	c.DeleteStatus(skillKey)
	c.SetCD(action.ActionSkill, 9*60)
	c.ConsumeSerpentsSubtlety(0, c.Base.Key.String()+"-skill-exit")
}

func (c *char) serpentsReduceTask(src int) {
	const tickInterval = .2
	c.Core.Tasks.Add(func() {
		if c.skillSrc != src {
			return
		}
		// reduce 1.4 point every 12f, which is 7 per second
		c.ReduceSerpentsSubtlety(c.Base.Key.String()+"skill", 1.4)
		if c.serpentsSubtlety == 0 && c.StatusIsActive(skillKey) {
			c.exitSkillState(src)
		}
		c.serpentsReduceTask(src)
	}, 60*tickInterval)
}

func (c *char) skillHold(p map[string]int) (action.Info, error) {
	duration := p["hold"]
	if duration < 10 {
		duration = 10
	}
	c.AddSerpentsSubtlety(c.Base.Key.String()+"-skill-hold", 45.0)
	c.c2OnSkill()

	c.absorbVoidRift()
	c.SetCDWithDelay(action.ActionSkill, 9*60, duration)

	return action.Info{
		Frames: func(next action.Action) int {
			return skillFrames[next] + duration
		},
		AnimationLength: skillFrames[action.InvalidAction] + duration,
		CanQueueAfter:   skillFrames[action.ActionDash] + duration, // earliest cancel is before skillHitmark
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
	c.AddStatus(particleICDKey, 15*60, false)

	count := 4.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Cryo, c.ParticleDelay)
}
