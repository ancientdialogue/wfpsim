package clorinde

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFrames   []int
	burstHitmarks = []int{104, 110, 116, 122, 128}
)

const (
	burstCD = 15 * 60
)

func init() {
	burstFrames = frames.InitAbilSlice(129)
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// add bol?
	c.ModifyHPDebtByRatio(burstBOL[c.TalentLvlBurst()])
	for _, v := range burstHitmarks {
		if c.Base.Cons >= 4 {
			c.c4()
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Burst",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeSlash,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       burstDamage[c.TalentLvlBurst()],
		}
		// TODO: what's the size of this??
		ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -1}, 11.2, 9)
		c.Core.QueueAttack(ai, ap, v, v)
	}

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(6)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}
