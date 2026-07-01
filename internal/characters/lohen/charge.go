package lohen

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

var chargeHitmarks = []int{20, 27}

func init() {
	chargeFrames = frames.InitAbilSlice(40)
	chargeFrames[action.ActionSkill] = chargeHitmarks[1]
	chargeFrames[action.ActionBurst] = chargeHitmarks[1]
	chargeFrames[action.ActionDash] = chargeHitmarks[1]
	chargeFrames[action.ActionJump] = chargeHitmarks[1]
	chargeFrames[action.ActionWalk] = 40
	chargeFrames[action.ActionSwap] = 40
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillChargeAttack()
	}

	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Charge",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagNormalAttack,
		ICDGroup:           attacks.ICDGroupLohenSkillAttack,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charge[c.TalentLvlAttack()],
		CanBeDefenseHalted: true,
	}
	ap := combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		info.Point{Y: 1.5},
		3.3,
		3,
	)
	for _, hitmark := range chargeHitmarks {
		c.Core.QueueAttack(ai, ap, hitmark, hitmark)
	}
	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[1],
		State:           action.ChargeAttackState,
	}, nil
}

func (c *char) skillChargeAttack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:         c.Index(),
		Abil:               "Charge (Masterstroke)",
		AttackTag:          attacks.AttackTagExtra,
		ICDTag:             attacks.ICDTagExtraAttack,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSlash,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skillCharge[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
		IgnoreInfusion:     true,
	}

	ap := combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		info.Point{Y: 1.5},
		3.3,
		3,
	)

	for _, hitmark := range chargeHitmarks {
		c.Core.QueueAttack(ai, ap, hitmark, hitmark, c.particleCB, c.joyCB, c.c2MakeCB())
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmarks[1],
		State:           action.ChargeAttackState,
	}, nil
}
