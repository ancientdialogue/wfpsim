package lohen

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames         []int
	skillPierceFrames   []int
	skillPierceHitmarks = []int{20, 20 + 3, 20 + 3 + 3, 20 + 3 + 3 + 3}
)

const (
	skillHitmark         = 19
	particleICDKey       = "lohen-particle-icd"
	skillKey             = "lohen-masterstroke"
	joyICDKey            = "lohen-joy-icd"
	joyMax               = 100
	defaultMaxSkillsUsed = 3
)

func init() {
	skillFrames = frames.InitAbilSlice(44)
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionSkill] = 22
	skillFrames[action.ActionBurst] = 19
	skillFrames[action.ActionDash] = 17
	skillFrames[action.ActionJump] = 18
	skillFrames[action.ActionSwap] = 17

	skillPierceFrames = frames.InitAbilSlice(42)
	skillPierceFrames[action.ActionAttack] = 28
	skillPierceFrames[action.ActionBurst] = 28
	skillPierceFrames[action.ActionDash] = 26
	skillPierceFrames[action.ActionJump] = 26
	skillPierceFrames[action.ActionWalk] = 32
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		// do nothing if previous char wasn't lohen
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		c.DeleteStatus(skillKey)
	}, "lohen-exit")
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillPierce(), nil
	}
	c.willToWin = 0
	c.joy = 0
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "Joy consumed (0)")
	c.maxSkillsUsed = defaultMaxSkillsUsed
	c.skillsUsed = 0
	// trigger 0 damage attack; matters because this stacks PJWS
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Strategem (0 dmg)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), skillHitmark, skillHitmark)

	c.AddStatus(skillKey, 13*60+skillHitmark, true)
	c.SetCDWithDelay(action.ActionSkill, 18*60, skillHitmark)

	c.skillBonusOnSkillMasterstroke()
	c.c4OnSkillMasterstroke()

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Cryo, c.ParticleDelay)
}

func (c *char) joyCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.skillsUsed >= c.maxSkillsUsed {
		return
	}

	if c.StatusIsActive(joyICDKey) {
		return
	}

	c.AddStatus(joyICDKey, 0.1*60, true)
	c.joy = min(c.joy+17, joyMax)
	c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), fmt.Sprintf("Joy gained (%v)", c.joy))
}

func (c *char) skillInit() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		if !c.StatusIsActive(skillKey) {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex == c.Index() {
			return
		}
		willToWinGen := 1.

		amt := args[2].(float64)
		if amt > c.Stat(attributes.BaseATK)*10 {
			willToWinGen = 20
		}

		willToWinGen += c.a1BonusWill(amt)

		c.gainWillToWin(willToWinGen, "skill")
	}, "lohen-will-to-win")
}

func (c *char) gainWillToWin(amt float64, src string) {
	amt *= c.c1WillToWinMult()
	c.willToWin = min(c.willToWin+amt, c.willToWinMax)

	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), fmt.Sprintf("Will to Win gained (%.2f)", c.willToWin)).
			Write("amt", amt).
			Write("abil", src)
	}
}

func (c *char) skillPierce() action.Info {
	// c6 duration extension happens immediately on skill use
	// to ensure that the skill doesn't expire during the animation,
	// causing the extension to fail
	c.c6ExtendOnSkill()

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Etched Into Bone and Soul",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       skillPierce[c.TalentLvlSkill()],
	}

	var snapshot info.Snapshot

	c.QueueCharTask(func() {
		snapshot = c.Snapshot(&ai)
		will := c.consumeWill()
		ai.Mult *= (1 + will*0.004)
		c.hexereiOnSkillBurst(will)
		c.c2OnSkillBurst()
		c.c6OnSkillBurst(will, &snapshot)
	}, skillPierceHitmarks[0]-1)

	ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 5, 60)
	for _, delay := range skillPierceHitmarks {
		c.QueueCharTask(func() { c.Core.QueueAttackWithSnap(ai, snapshot, ap, 0) }, delay)
	}

	return action.Info{
		Frames:          func(next action.Action) int { return skillPierceFrames[next] },
		AnimationLength: skillPierceFrames[action.InvalidAction],
		CanQueueAfter:   skillPierceFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) consumeWill() float64 {
	if c.StatusIsActive(skillKey) {
		will := c.willToWin
		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "Will to Win Consumed").
				Write("amt", will).
				Write("mult", 1+will*0.004)
		}

		c.willToWin = 0

		return will
	}
	c.willToWin = 0
	return 0
}
