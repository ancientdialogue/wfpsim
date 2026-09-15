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

const burstKey = "alyosha-burst"

func init() {
	burstFrames = frames.InitAbilSlice(55) // Q -> E
	burstFrames[action.ActionAttack] = 54  // Q -> N1
	burstFrames[action.ActionSkill] = 53   // Q -> E
	burstFrames[action.ActionWalk] = 54    // Q -> W
	burstFrames[action.ActionSwap] = 52    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	src := c.Core.F
	c.burstSrc = src
	c.Core.Tasks.Add(func() { c.burstTicker(src) }, 79)
	c.AddStatus(burstKey, 14*60+c.c2BurstDur(), true)

	c.SetCD(action.ActionBurst, 18*60)
	c.ConsumeEnergy(7)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
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
		Abil:       "Burst",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupAlyoshaBurst,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0, c.c2MakeBurstCB())

	aiDog := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Tugarin",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupAlyoshaBurst,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       burstTick[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(aiDog, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 39, 39, c.c2MakeBurstCB())
	c.a1OnBurstTick()
	c.c4OnBurstTick()
	c.Core.Tasks.Add(func() { c.burstTicker(src) }, 119)
}
