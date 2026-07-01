package illuga

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstHitmark = 31 // Initial Hit
	burstKey     = "illuga-burst"
	initalStacks = 21
	maxGained    = 15
)

func init() {
	burstFrames = frames.InitAbilSlice(56) // Q -> E
	burstFrames[action.ActionAttack] = 53  // Q -> N1
	burstFrames[action.ActionDash] = 42    // Q -> D
	burstFrames[action.ActionJump] = 43    // Q -> J
	burstFrames[action.ActionSwap] = 55    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstStacks = 0
	c.stacksGained = 0
	c.QueueCharTask(func() {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Burst",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			PoiseDMG:   40,
			Element:    attributes.Geo,
			Durability: 25,
			Mult:       burstDef[c.TalentLvlBurst()],
			FlatDmg:    burstEM[c.TalentLvlBurst()] * c.Stat(attributes.EM),
			UseDef:     true,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, 0)

		c.AddStatus(burstKey, 20*60, true)
		c.burstStacks = initalStacks
		c.gainBurstStacks(c.Core.Constructs.Count() * 5)
	}, burstHitmark)

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(7)
	c.a1OnSkillBurst()

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstInit() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		ae := args[1].(*info.AttackEvent)
		if ae.Info.Element != attributes.Geo {
			return
		}

		switch ae.Info.AttackTag {
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		case attacks.AttackTagExtra:
		case attacks.AttackTagNormal:
		case attacks.AttackTagPlunge:
		case attacks.AttackTagDirectLunarCrystallize:
		default:
			return
		}

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if !c.StatusIsActive(burstKey) {
			return
		}

		if c.burstStacks == 0 {
			return
		}

		c.useStack()

		ratio := geoBonus[c.TalentLvlBurst()]
		isLCr := ae.Info.AttackTag == attacks.AttackTagDirectLunarCrystallize
		if isLCr {
			ratio = lcrBonus[c.TalentLvlBurst()]
		}

		ratio += c.a4BonusScaling(isLCr)

		em := c.Stat(attributes.EM)
		amt := ratio * em
		c.Core.Log.NewEvent("Illuga burst proc dmg add", glog.LogPreDamageMod, ae.Info.ActorIndex).
			Write("em", em).
			Write("ratio", ratio).
			Write("addition", amt)

		ae.Info.FlatDmg += amt
	}, "illuga-burst-hook")
}

func (c *char) gainBurstStacks(amt int) {
	amt = min(amt, maxGained-c.stacksGained)
	c.burstStacks += amt
	c.stacksGained += amt
}

func (c *char) useStack() {
	c.burstStacks -= 1
	c.c2OnStackConsume()
}
