package pyro

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

const chargeTruemoonICDKey = "travelerpyro-special-ca-icd"

var (
	chargeFrames   [][]int
	chargeHitmarks = [][]int{{9, 20}, {14, 25}}
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
	conversion := c.chargeAttackTruemoon()
	if conversion == nil {
		conversion = c.c6Conversion
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
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
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames[c.gender]),
		AnimationLength: chargeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[c.gender][len(chargeHitmarks[c.gender])-1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *Traveler) chargeAttackTruemoon() func(*info.AttackInfo) {
	if !c.trueMoonBuff {
		return nil
	}

	if c.trueMoonStacks < 2 {
		return nil
	}

	if c.StatusIsActive(chargeTruemoonICDKey) {
		return nil
	}
	c.AddStatus(chargeTruemoonICDKey, 15*60, true)
	c.trueMoonStacks = 0
	return func(ai *info.AttackInfo) {
		ai.Element = attributes.Pyro
		ai.IgnoreInfusion = true
		ai.ICDTag = attacks.ICDTagTravelerEnchancedCA
		ai.AdditionalTags = []attacks.AdditionalTag{attacks.AdditionalTagNightsoul}
		ai.Mult += 2
		ai.Abil += " (Inferno)"
	}
}

func (c *Traveler) trueMoonInit() {
	if !c.trueMoonBuff {
		return
	}

	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		c.trueMoonStacks = min(c.trueMoonStacks+1, 2)
	}, "travelerpyro-truemoon")
}
