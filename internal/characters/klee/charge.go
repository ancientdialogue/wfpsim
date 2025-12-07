package klee

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	chargeHitmark  = 76
	chargeSnapshot = 29 + 32
)

var (
	boombadgeMult = []float64{1, 1.15, 1.3, 1.5}
	chargeFrames  []int
)

func init() {
	chargeFrames = frames.InitAbilSlice(113)
	chargeFrames[action.ActionAttack] = 59
	chargeFrames[action.ActionCharge] = 59
	chargeFrames[action.ActionSkill] = 59
	chargeFrames[action.ActionBurst] = 59
	chargeFrames[action.ActionDash] = 31
	chargeFrames[action.ActionJump] = 30
	chargeFrames[action.ActionSwap] = 104
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	windup := 0
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState:
		if c.NormalCounter == 1 || c.NormalCounter == 2 {
			windup = 14
		}
	case action.SkillState:
		windup = 14
	}

	c.Core.Tasks.Add(func() {
		ai := c.getChargeAttackInfo()
		snap := c.applySpark(&ai)
		c.Core.QueueAttackWithSnap(
			ai,
			snap,
			combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 3),
			(chargeHitmark-chargeSnapshot)+travel,
			c.makeA4CB(),
		)
	}, chargeSnapshot-windup)

	c.c1(chargeHitmark - windup + travel)

	return action.Info{
		Frames:          func(next action.Action) int { return chargeFrames[next] - windup },
		AnimationLength: chargeFrames[action.InvalidAction] - windup,
		CanQueueAfter:   chargeFrames[action.ActionJump] - windup, // earliest cancel
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) applySpark(ai *info.AttackInfo) info.Snapshot {
	snap := c.Snapshot(ai)
	if c.StatusIsActive(a1SparkKey) {
		ai.Abil = "Boom-Boom Strike"
		snap.Stats[attributes.DmgP] += .50
		// Hexerei bonus (C6):
		// When Klee uses an Explosive Spark, there is a 50% chance it will not be consumed.
		previous := c.a1CurrentStack
		if c.Base.Cons < 6 || c.IsHexerei && c.Core.Player.GetHexereiCount() > 1 && c.Core.Rand.Float64() < 0.5 {
			c.a1CurrentStack--
		}
		c.Core.Log.NewEvent("consuming spark", glog.LogCharacterEvent, c.Index()).
			Write("previous spark stacks", previous).
			Write("new spark stacks", c.a1CurrentStack).
			Write("boombadge mult", c.getBoomBadgeMult())
		if c.a1CurrentStack == 0 {
			c.DeleteStatus(a1SparkKey)
		}
	}
	return snap
}

func (c *char) getBoomBadgeStacks() int {
	count := 0
	if c.StatusIsActive(boomBadgeNormalKey) {
		count++
	}
	if c.StatusIsActive(boomBadgeSkillKey) {
		count++
	}
	if c.StatusIsActive(boomBadgeBurstKey) {
		count++
	}
	return count
}

// TODO: Is boombadge applied on snapshot?
func (c *char) getBoomBadgeMult() float64 {
	if !c.StatusIsActive(a1SparkKey) {
		return 1.0
	}
	return boombadgeMult[c.getBoomBadgeStacks()]
}

func (c *char) getChargeAttackInfo() info.AttackInfo {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		PoiseDMG:   180,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()] * c.getBoomBadgeMult(),
	}
	return ai
}
