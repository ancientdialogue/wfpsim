package fischl

import (
	"github.com/genshinsim/gcsim/internal/template/minazuki"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

const c6HexereiKey = "fischl-c6-hexerei"

func (c *char) c6Init() error {
	if c.Base.Cons < 6 {
		return nil
	}

	w, err := minazuki.New(
		minazuki.WithMandatory(keys.Fischl, "fischl c6", ozActiveKey, "", 60, c.c6Wave, c.Core),
		minazuki.WithTickOnActive(true),
		minazuki.WithAnimationDelayCheck(info.AnimationYelanN0StartDelay, func() bool {
			return c.Core.Player.ActiveChar().NormalCounter == 1
		}),
	)
	if err != nil {
		return err
	}
	c.c6Watcher = w
	return nil
}

func (c *char) c6Wave() {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Evernight Raven (C6)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagElementalArt,
		ICDGroup:   attacks.ICDGroupFischl,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       0.3,
	}

	// C6 uses Oz Snapshot
	c.Core.QueueAttackWithSnap(
		ai,
		c.ozSnapshot.Snapshot,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			info.Point{Y: -1},
			0.1,
			1,
		),
		c.ozTravel,
	)

	if c.IsHexerei {
		c.AddStatus(c6HexereiKey, 10*60, true)
	}
}

func (c *char) c6HexereiBonus() float64 {
	if c.Base.Cons < 6 {
		return 1.0
	}

	if !c.StatusIsActive(c6HexereiKey) {
		return 1.0
	}
	return 2.0
}
