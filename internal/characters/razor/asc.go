package razor

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Decreases Claw and Thunder's CD by 18%.
func (c *char) a1CDReduction(cd int) int {
	if c.Base.Ascension < 1 {
		return cd
	}
	return int(float64(cd) * 0.82)
}

// Using Lightning Fang resets the CD of Claw and Thunder.
func (c *char) a1CDReset() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// When Razor's Energy is below 50%, increases Energy Recharge by 30%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.a4Bonus = make([]float64, attributes.EndStatType)
	c.a4Bonus[attributes.ER] = 0.3
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("razor-a4", -1),
		AffectedStat: attributes.ER,
		Amount: func() ([]float64, bool) {
			if c.Energy/c.EnergyMax >= 0.5 {
				return nil, false
			}
			return c.a4Bonus, true
		},
	})
}

func (c *char) magicWolfMult() float64 {
	// It seems like he gets the Q buff regardless if he is a magic character?
	return 0.7
}

const (
	magicSigilIcdKey = "razor-magic-icd"
	magicSigilKey    = "razor-magic"
)

func (c *char) magicOnSigilOverflow() {
	if !c.IsMagic {
		return
	}

	if c.getMagicCount() < 2 {
		return
	}

	if !c.StatusIsActive(burstBuffKey) {
		return
	}

	if c.StatusIsActive(magicSigilIcdKey) {
		return
	}
	c.AddStatus(magicSigilIcdKey, 60, true)

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Magic: Secret Rite (Razor)",
		AttackTag:  attacks.AttackTagElementalArt, // TODO: it has another tag?
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       1.5,
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1.5),
		6,
		6,
	)
	c.AddEnergy(magicSigilKey, 7)
	c.c6OnSiglConsume()
}

func (c *char) getMagicCount() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		if c.IsMagic {
			count += 1
		}
	}
	return count
}
