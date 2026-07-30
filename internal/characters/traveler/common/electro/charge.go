package electro

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

const (
	chargeTruemoonICDKey = "travelergeo-special-ca-icd"
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
	amuletDelay := max(
		// make it so that it can't be faster than 3+1.79s
		p["amulet_delay"],
		// ~1.79s
		3*60+107)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNormalAttack,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Physical,
		Durability: 25,
	}
	conversion := c.chargeAttackTruemoon(amuletDelay)

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

func (c *Traveler) chargeAttackTruemoon(amuletDelay int) func(*info.AttackInfo) {
	if !c.trueMoonBuff {
		return func(*info.AttackInfo) {}
	}

	if c.trueMoonStacks < 3 {
		return func(*info.AttackInfo) {}
	}

	if c.StatusIsActive(chargeTruemoonICDKey) {
		return func(*info.AttackInfo) {}
	}

	c.AddStatus(chargeTruemoonICDKey, 15*60, true)
	c.trueMoonStacks = 0

	lightning := func() {
		ai := info.AttackInfo{
			ActorIndex:     c.Index(),
			AttackTag:      attacks.AttackTagExtra,
			ICDTag:         attacks.ICDTagNone,
			ICDGroup:       attacks.ICDGroupDefault,
			StrikeType:     attacks.StrikeTypeSlash,
			Element:        attributes.Electro,
			Durability:     25,
			Mult:           2,
			IgnoreInfusion: true,
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), info.Point{Y: -1}, 1.5)
		c.Core.QueueAttack(ai, ap, 0, 0)

		// 1 amulet per attack
		c.abundanceAmulets++
		c.SetTag("generated", c.abundanceAmulets)

		c.Core.Log.NewEvent("travelerelectro abundance amulet generated", glog.LogCharacterEvent, c.Index()).
			Write("amulets", c.abundanceAmulets)
	}
	c.Core.Tasks.Add(lightning, chargeHitmarks[c.gender][len(chargeHitmarks[c.gender])-1]+3*60)

	c.Core.Tasks.Add(func() {
		active := c.Core.Player.ActiveChar()
		c.collectAmulets(active)
	}, amuletDelay)

	return func(ai *info.AttackInfo) {
		ai.Element = attributes.Geo
		ai.IgnoreInfusion = true
		ai.ICDTag = attacks.ICDTagTravelerEnchancedCA
		ai.Mult += 1
		ai.Abil += " (Detonate)"
	}
}

func (c *Traveler) trueMoonInit() {
	if !c.trueMoonBuff {
		return
	}

	gainStacks := func(args ...any) {
		c.trueMoonStacks = min(c.trueMoonStacks+1, 3)
	}
	c.Core.Events.Subscribe(event.OnBurst, gainStacks, "travelergeo-truemoon")
}
