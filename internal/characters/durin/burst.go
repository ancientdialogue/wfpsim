package durin

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames      []int
	burstInitHitmark = []int{110, 110 + 10, 110 + 10 + 10} // Initial Hit
)

const (
	burstTicks          = 16
	burstInterval       = 75
	burstFirstTickDelay = 140
	burstCD             = 18 * 60
	burstKeyWhite       = "durin-burst-white"
	burstKeyBlack       = "durin-burst-black"
)

func init() {
	burstFrames = frames.InitAbilSlice(115) // E -> D/J
	burstFrames[action.ActionAttack] = 115
	burstFrames[action.ActionBurst] = 115
	burstFrames[action.ActionWalk] = 115
	burstFrames[action.ActionSwap] = 115
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(blackKey) {
		return c.burstBlack()
	}

	return c.burstWhite()
}

func (c *char) burstWhite() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Lustrous Light: Birthed by Dusk",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		FlatDmg:    c.a4Dmg(),
	}
	for i, mult := range burstWhiteInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 5), burstInitHitmark[i], burstInitHitmark[i])
	}

	c.burstSrc = c.Core.F
	for i := 0.0; i < burstTicks; i++ {
		c.QueueCharTask(c.burstTickWhite(c.burstSrc), burstFirstTickDelay+ceil(burstInterval*i))
	}
	c.DeleteStatus(burstKeyBlack)
	c.AddStatus(burstKeyWhite, burstFirstTickDelay+ceil((burstTicks-1)*burstInterval), false)

	c.SetCDWithDelay(action.ActionBurst, burstCD, 22)
	c.ConsumeEnergy(10)
	c.a1OnBurst(true)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTickWhite(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Lustrous Light: Searing Flame",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDurin,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       burstWhiteDoT[c.TalentLvlBurst()],
			FlatDmg:    c.a4Dmg(),
		}

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5), 0, 0)
	}
}

func (c *char) burstBlack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Dark Decay: Devoured by Dawn",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		FlatDmg:    c.a4Dmg(),
	}
	for i, mult := range burstBlackInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 5), burstInitHitmark[i], burstInitHitmark[i])
	}

	c.burstSrc = c.Core.F
	for i := 0.0; i < burstTicks; i++ {
		c.QueueCharTask(c.burstTickBlack(c.burstSrc), burstFirstTickDelay+ceil(burstInterval*i))
	}
	c.DeleteStatus(burstKeyWhite)
	c.AddStatus(burstKeyBlack, burstFirstTickDelay+ceil((burstTicks-1)*burstInterval), false)

	c.SetCDWithDelay(action.ActionBurst, burstCD, 22)
	c.ConsumeEnergy(10)
	c.a1OnBurst(false)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTickBlack(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Dark Decay: Abyssal Flame",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDurin,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       burstBlackDoT[c.TalentLvlBurst()],
			FlatDmg:    c.a4Dmg(),
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0)
	}
}
