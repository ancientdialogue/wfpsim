package lohen

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames   []int
	burstHitmarks = []int{101, 101 + 21, 101 + 21 + 6, 101 + 21 + 6 + 9, 101 + 21 + 6 + 9 + 25}
)

func init() {
	burstFrames = frames.InitAbilSlice(130) // Q -> W
	burstFrames[action.ActionAttack] = 129  // Q -> N1
	burstFrames[action.ActionSkill] = 129   // Q -> E
	burstFrames[action.ActionSwap] = 129    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		c.ExtendStatus(skillKey, 1.65*60)
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Manifest Judgement",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Cryo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -5},
		14,
		12,
	)
	var snapshot info.Snapshot
	c.QueueCharTask(func() {
		snapshot = c.Snapshot(&ai)
		c.c4OnBurst()
		will := c.consumeWill()
		ai.Mult *= (1 + will*0.004)
		c.hexereiOnSkillBurst(will)
		c.c2OnSkillBurst()
		c.c6OnSkillBurst(will, &snapshot)
	}, burstHitmarks[0]-1)

	for _, delay := range burstHitmarks {
		c.QueueCharTask(func() { c.Core.QueueAttackWithSnap(ai, snapshot, ap, 0) }, delay)
	}

	c.ConsumeEnergy(7)
	c.QueueCharTask(func() { c.c4OnBurstEnergyConsume() }, 7)

	c.SetCDWithDelay(action.ActionBurst, 15*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
