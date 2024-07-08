package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Cleardew Cologne",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       6.0,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 4.5}, 3)
	c.Core.QueueAttack(ai, ap, 0, 30)
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = min(c.getTotalAtk()/1000*0.15, 0.36)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("emilie-a4", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if e, ok := t.(*enemy.Enemy); !(ok && e.IsBurning()) {
				return nil, false
			}
			return m, true
		},
	})
}
