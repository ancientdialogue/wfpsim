package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

const (
	burstKey      = "mavuika-burst"
	burstDuration = 7 * 60
	burstHitmark  = 110
)

var (
	burstFrames []int
)

func (c *char) nightsoulConsumptionMul() float64 {
	if c.StatusIsActive(burstKey) {
		return 0.0
	}
	return 1.0
}

func init() {
	burstFrames = frames.InitAbilSlice(120) // Q -> N1
	burstFrames[action.ActionSkill] = 120
	burstFrames[action.ActionDash] = 120
	burstFrames[action.ActionJump] = 120
	burstFrames[action.ActionSwap] = burstHitmark
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstStacks = c.fightingSpirit
	c.fightingSpirit = 0
	c.armamentState = bike
	c.enterNightsoulOrRegenerate(10)

	c.QueueCharTask(func() {
		c.AddStatus(burstKey, burstDuration, true)
		c.a4()

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Sunfell Slice",
			AttackTag:      attacks.AttackTagElementalBurst,
			ICDTag:         attacks.ICDTagNone,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeBlunt,
			PoiseDMG:       150,
			Element:        attributes.Pyro,
			Durability:     25,
			Mult:           burst[c.TalentLvlBurst()],
			FlatDmg:        c.burstBuffSunfell(),
		}
		ap := combat.NewCircleHitOnTarget(
			c.Core.Combat.Player(),
			geometry.Point{Y: 1.0},
			6,
		)
		c.Core.QueueAttack(ai, ap, 0, 0)
	}, burstHitmark)

	c.SetCDWithDelay(action.ActionBurst, 18*60, 0)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstBuffCA() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstCABonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) burstBuffNA() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstNABonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) burstBuffSunfell() float64 {
	if !c.StatusIsActive(burstKey) {
		return 0.0
	}
	return c.burstStacks * burstQBonus[c.TalentLvlBurst()] * c.TotalAtk()
}

func (c *char) gainFightingSpirit(val float64) {
	c.fightingSpirit += val * c.c1FightingSpiritEff()
	if c.fightingSpirit > 200 {
		c.fightingSpirit = 200
	}
	c.c1OnFightingSpirit()
}

func (c *char) burstInit() {
	c.fightingSpirit = 200
	c.Core.Events.Subscribe(event.OnNightsoulConsume, func(args ...interface{}) bool {
		amount := args[1].(float64)
		if amount < 0.0000001 {
			return false
		}
		c.gainFightingSpirit(amount)
		return false
	}, "mavuika-fighting-spirit")
}