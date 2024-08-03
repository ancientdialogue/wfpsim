package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
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

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("emilie-a4", -1),
		Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				return nil, false
			}
			if !e.IsBurning() {
				return nil, false
			}
			m[attributes.DmgP] = 0.00015 * c.TotalAtk()
			return m, true
		},
	})
}
