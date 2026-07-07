package cryo

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	chargeFrames   [][]int
	chargeHitmarks = [][]int{{9, 20}, {14, 25}}
)

const (
	chargeFreezingIceICDKey = "travelercryo-special-ca-icd"
	trueMoonStackICDKEy     = "travelercryo-icepoint-icd"
)

func init() {
	chargeFrames = make([][]int, 2)
	// Male
	chargeFrames[0] = frames.InitAbilSlice(55)                                       // CA -> N1
	chargeFrames[0][action.ActionSkill] = 37                                         // CA -> E
	chargeFrames[0][action.ActionBurst] = 36                                         // CA -> Q
	chargeFrames[0][action.ActionDash] = chargeHitmarks[0][len(chargeHitmarks[0])-1] // CA -> D
	chargeFrames[0][action.ActionJump] = chargeHitmarks[0][len(chargeHitmarks[0])-1] // CA -> J
	chargeFrames[0][action.ActionSwap] = 44                                          // CA -> Swap

	// Female
	chargeFrames[1] = frames.InitAbilSlice(58)                                       // CA -> N1
	chargeFrames[1][action.ActionSkill] = 34                                         // CA -> E
	chargeFrames[1][action.ActionBurst] = 35                                         // CA -> Q
	chargeFrames[1][action.ActionDash] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> D
	chargeFrames[1][action.ActionJump] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> J
	chargeFrames[1][action.ActionSwap] = chargeHitmarks[1][len(chargeHitmarks[1])-1] // CA -> Swap
}

func (c *Traveler) ChargeAttack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}

	conversion := c.chargeAttackFreezingIce()
	if conversion == nil {
		conversion = c.a1Conversion
	}

	for i, mult := range charge[c.gender] {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		conversion(&ai)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.2),
			chargeHitmarks[c.gender][i],
			chargeHitmarks[c.gender][i],
			c.naCaPlungeCB,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames[c.gender]),
		AnimationLength: chargeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[c.gender][len(chargeHitmarks[c.gender])-1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *Traveler) chargeAttackFreezingIce() func(*info.AttackInfo) {
	if !c.trueMoonBuff {
		return nil
	}

	if c.trueMoonStacks < 3 {
		return nil
	}

	if !c.StatusIsActive(skillKey) {
		return nil
	}

	if c.StatusIsActive(chargeFreezingIceICDKey) {
		return nil
	}
	c.AddStatus(chargeFreezingIceICDKey, 15*60, true)
	c.trueMoonStacks = 0
	c.flostglowStacks = min(c.flostglowStacks+2, skillStacksMax)
	return func(ai *info.AttackInfo) {
		ai.Element = attributes.Cryo
		ai.IgnoreInfusion = true
		ai.Mult += 1.4

		switch c.getRadiance() {
		case radianceStellarConduct:
			ai.Abil += stellarConductText
			ai.AttackTag = attacks.AttackTagDirectStellarConduct
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
		case radianceStellarSwirl:
			ai.Abil += stellarSwirlText
			ai.AttackTag = attacks.AttackTagDirectStellarSwirl
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
		}
	}
}

func (c *Traveler) trueMoonInit() {
	if !c.trueMoonBuff {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		case attacks.AttackTagReactionStellarSwirl:
		default:
			return
		}

		if c.StatusIsActive(trueMoonStackICDKEy) {
			return
		}

		c.AddStatus(trueMoonStackICDKEy, 2*60, true)

		c.trueMoonStacks = min(c.trueMoonStacks+1, 3)
	}, "cryo-mc-truemoon")
}
