package columbina

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(112)
	burstFrames[action.ActionSwap] = 111
}

const (
	burstBuffKey = "columbina-q-buff"
	burstDur     = 20 * 60
)

func (c *char) Burst(p map[string]int) (action.Info, error) {
	c.burstArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, 20)
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Burst",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		UseHP:      true,
		Mult:       burst[c.TalentLvlBurst()],
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6)
	c.Core.QueueAttack(ai, ap, 105, 105)

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(burstBuffKey, burstDur+105),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				switch ai.AttackTag {
				case attacks.AttackTagDirectLunarCharged:
				case attacks.AttackTagReactionLunarCharge:
				default:
					return 0, false
				}

				if !c.Core.Combat.Player().IsWithinArea(c.burstArea) {
					return 0, false
				}
				if c.Core.Combat.Debug {
					c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), "Adding columbina burst react bonus")
				}
				return burstBuff[c.TalentLvlBurst()], false
			},
		})
	}

	c.ConsumeEnergy(5)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
