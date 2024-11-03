package chasca

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	if len(c.shadowhuntShells) >= 3 {
		uniqueCount := len(c.uniqueConversionElements)
		chance := 0.0
		if c.Base.Cons > 1 {
			chance = 0.333
		}
		switch uniqueCount {
		case 1:
			chance += 0.333
		case 2:
			chance += 0.667
		default:
			chance += 1.0
		}

		if c.Core.Rand.Float64() < chance {
			if len(c.conversionElements) > 0 {
				randomIndex := c.Core.Rand.Intn(len(c.conversionElements))
				c.shadowhuntShells[2] = c.conversionElements[randomIndex]
				if c.Base.Cons > 1 {
					c1RandomIndex := c.Core.Rand.Intn(len(c.conversionElements))
					c.shadowhuntShells[1] = c.conversionElements[c1RandomIndex]
				}
			} else {
				c.Core.Log.NewEvent("chasca a1: conversionElements is empty", glog.LogWarnings, -1)
			}
		}
	}
}

func (c *char) a1Amount() float64 {
	a1Boost := 0.0
	if c.Base.Ascension < 1 {
		return a1Boost
	}
	uniqueCount := len(c.uniqueConversionElements)
	switch uniqueCount {
	case 1:
		a1Boost = 0.15
	case 2:
		a1Boost = 0.35
	default:
		a1Boost = 0.65
	}
	mDmg := make([]float64, attributes.EndStatType)
	mDmg[attributes.DmgP] = a1Boost
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("chasca-a1-dmg-bonus", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.Abil != "Shining Shadowhunt Shell" {
				return nil, false
			}
			return mDmg, true
		},
	})
	return a1Boost
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.Core.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
		element := attributes.Anemo
		for _, char := range c.Core.Player.Chars() {
			switch char.Base.Element {
			case attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro:
				element = char.Base.Element
				break
			}
		}
		if element == attributes.Anemo {
			c.a4Dmg = 1.5 * skillShadowhunt[c.TalentLvlSkill()]
		} else {
			c.a4Dmg = 1.5 * skillShiningShadowhunt[c.TalentLvlSkill()]
		}

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Burning Shadowhunt Shot",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    element,
			Durability: 25,
			Mult:       c.a4Dmg,
		}

		target := c.Core.Combat.ClosestEnemyWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 15), nil)
		if target != nil {
			ap := combat.NewCircleHitOnTarget(target, nil, 1)
			c.Core.QueueAttack(ai, ap, 0, 1)
		}

		return false
	}, "chasca-a4")
}
