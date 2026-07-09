package alyosha

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
	burstHitmark = 27 // Initial Hit
	burstKey     = "alyosha-burst"
)

func init() {
	burstFrames = frames.InitAbilSlice(56) // Q -> E
	burstFrames[action.ActionAttack] = 53  // Q -> N1
	burstFrames[action.ActionDash] = 42    // Q -> D
	burstFrames[action.ActionJump] = 43    // Q -> J
	burstFrames[action.ActionSwap] = 55    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Burst",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       burst[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0)

		c.AddStatus(burstKey, 14*60+c.c2BurstDur(), true)

		src := c.Core.F
		c.burstSrc = src
		c.Core.Tasks.Add(func() { c.burstTicker(src) }, 48)
	}, burstHitmark)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTicker(src int) {
	if c.burstSrc != src {
		return
	}

	if !c.StatusIsActive(burstKey) {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Tugarin",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupAlyoshaBurst,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0, c.c2MakeBurstCB())
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 5, 5, c.c2MakeBurstCB())
	c.a1OnBurstTick()
	c.c4OnBurstTick()
	c.Core.Tasks.Add(func() { c.burstTicker(src) }, 119)
}
