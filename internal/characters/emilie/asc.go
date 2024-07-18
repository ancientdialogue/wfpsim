package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	if c.lumidouceScents < 2 {
		return
	}

	if c.lumidouceScents%2 > 0 {
		return
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cleardew Cologne",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       6.0,
	}
	ap := combat.NewCircleHit(
		c.lumidoucePos,
		c.Core.Combat.PrimaryTarget(),
		nil,
		5,
	)
	c.Core.QueueAttack(ai, ap, 5, 5, c.particleCB)
}
