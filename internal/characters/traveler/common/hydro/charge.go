package hydro

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	chargeTruemoonICDKey = "travelerhydro-special-ca-icd"
	trueMoonStackICDKey  = "travelerhydro-waters-icd"
)

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
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}
	conversion, cb := c.chargeAttackTruemoon()

	for i, mult := range charge[c.gender] {
		ai.Mult = mult[c.TalentLvlAttack()]
		ai.Abil = fmt.Sprintf("Charge %v", i)
		conversion(&ai)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.2),
			chargeHitmarks[c.gender][i],
			chargeHitmarks[c.gender][i],
			cb,
		)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames[c.gender]),
		AnimationLength: chargeFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[c.gender][len(chargeHitmarks[c.gender])-1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *Traveler) chargeAttackTruemoon() (func(*info.AttackInfo), info.AttackCBFunc) {
	if !c.trueMoonBuff {
		return func(*info.AttackInfo) {}, nil
	}

	if c.trueMoonStacks < 3 {
		return func(*info.AttackInfo) {}, nil
	}

	if c.StatusIsActive(chargeTruemoonICDKey) {
		return func(*info.AttackInfo) {}, nil
	}

	c.AddStatus(chargeTruemoonICDKey, 15*60, true)
	c.trueMoonStacks = 0

	bonusMV := 0.0
	var healingCB info.AttackCBFunc
	if c.CurrentHPRatio() >= 0.5 {
		c.Core.Player.Drain(info.DrainInfo{
			ActorIndex: c.Index(),
			Abil:       "Charged Attack: Tidebound",
			Amount:     c.MaxHP() * 0.1,
		})

		bonusMV = 1
	} else {
		done := false
		healingCB = func(a info.AttackCB) {
			if a.Target.Type() != info.TargettableEnemy {
				return
			}

			if done {
				return
			}
			done = true
			c.Core.Player.Heal(info.HealInfo{
				Caller:  c.Index(),
				Target:  c.Index(),
				Message: "Charged Attack: Tidebound",
				Src:     c.MaxHP() * 0.25,
				Bonus:   c.Stat(attributes.Heal),
			})
		}
	}

	return func(ai *info.AttackInfo) {
		ai.Element = attributes.Hydro
		ai.IgnoreInfusion = true
		ai.ICDTag = attacks.ICDTagTravelerEnchancedCA
		ai.Mult += 1.0 + bonusMV
		ai.Abil += " (Tidebound)"
	}, healingCB
}

func (c *Traveler) trueMoonInit() {
	if !c.trueMoonBuff {
		return
	}

	totalHpChange := 0.0

	gainStacks := func(amt float64) {
		totalHpChange += amt
		if totalHpChange < 0.05 {
			return
		}

		if c.StatusIsActive(trueMoonStackICDKey) {
			return
		}

		c.AddStatus(trueMoonStackICDKey, 4*60, true)
		c.trueMoonStacks = min(c.trueMoonStacks+1, 3)
	}

	c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(args ...any) {
		di := args[0].(*info.DrainInfo)

		if di.Amount <= 0 {
			return
		}

		char := c.Core.Player.ByIndex(di.ActorIndex)
		amt := di.Amount / char.MaxHP()
		gainStacks(amt)
	}, "travelerhydro-truemoon")

	c.Core.Events.Subscribe(event.OnHeal, func(args ...any) {
		target := args[1].(int)
		amount := args[2].(float64)
		overheal := args[3].(float64)

		if amount <= 0 {
			return
		}

		if math.Abs(amount-overheal) <= 1e-9 {
			return
		}

		char := c.Core.Player.ByIndex(target)
		amt := (amount - overheal) / char.MaxHP()

		gainStacks(amt)
	}, "travelerhydro-truemoon")
}
