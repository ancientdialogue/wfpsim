package clorinde

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1_icd              int     = 1 * 60
	c1_atkP             float64 = 0.3
	c1_icdKey                   = "clorinde-C1-IcdKey"
	c2_a1FlatDmg        float64 = 2700
	c2_a1PercentBuff    float64 = 0.3
	c6_icdKey                   = "clorinde-c6-icd"
	c6_mitigate                 = 0.8
	c6_glimbrightIcdKey         = "glimbrightIcdKey"
	c6_glimbrightAtkP           = 2
)

var c1_hitMarks = []int{1, 1} // TODO hitmark for each c1 hit

// While Hunt the Dark's Night Watch state is active,
// when Electro DMG from Clorinde's Normal Attacks hit opponents,
// they will trigger 2 coordinated attack from a Nightwatch Shade
// summoned near the hit opponent,
// each dealing 30% of Clorinde's ATK as Electro DMG.
// This effect can occur once every 1s.
// DMG dealt this way is considered Normal Attack DMG.

func (c *char) c1() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		if !c.StatusIsActive(skillStateKey) {
			return false
		}
		if c.StatusIsActive(c1_icdKey) {
			return false
		}
		c.AddStatus(c1_icdKey, c1_icd, false)
		c1_ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Nightwatch Shade (C1)",
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       c1_atkP,
		}
		target := args[0].(combat.Target)
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagNormal {
			for _, hitmark := range c1_hitMarks {
				c.Core.QueueAttack(
					c1_ai,
					combat.NewCircleHitOnTarget(target, nil, 4),
					hitmark,
					hitmark,
				)
			}
		}
		return false
	}, "clorinde-c1")
}

// When Last Lightfall deals DMG to opponent(s),
// DMG dealt is increased based on Clorinde's Bond of Life percentage.
// Every 1% of her current Bond of Life will increase Last Lightfall DMG by 2%.
// The maximum Last Lightfall DMG increase achievable this way is 200%.

func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = min(float64(c.currentHPDebtRatio())*100*0.02, 2)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("clorinde-c4-burst-bonus", 130), //TODO: bol snapshot frame
		AffectedStat: attributes.DmgP,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// For 12s after Hunt the Dark is used,
// Clorinde's CRIT Rate will be increased by 10%,
// and her CRIT DMG by 70%
func (c *char) c6skill() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.StatusIsActive(skillStateKey) {
		return
	}
	if c.StatusIsActive(c6_icdKey) {
		return
	}
	c.AddStatus(c6_icdKey, 12*60, true)

	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.1
	m[attributes.CD] = 0.7
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("clorinde-c6-bonus-cr", 12*60),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("clorinde-c6-bonus-cd", 12*60),
		AffectedStat: attributes.CD,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}

// Additionally, while Night Watch is active,
// a Glimbright Shade will appear under specific circumstances,
// decreasing DMG dealt to Clorinde by 80% for 1s
// and increasing her interruption resistance;
// it will also attack opponents,
// dealing 200% of Clorinde's ATK as Electro DMG.
// DMG dealt this way is considered Normal Attack DMG.
// The Glimbright Shade will appear under the following circumstances:
// · When Clorinde is about to be attacked by an opponent. TODO: currently not implemented
// · When Clorinde uses Impale the Night: Pact.
// 1 Glimbright Shade can be summoned in the aforementioned ways every 1s.
// 6 Shades can be summoned per single Night Watch duration.

func (c *char) c6() {
	if c.StatusIsActive(c6_glimbrightIcdKey) {
		return
	}

	c.AddStatus(c6_glimbrightIcdKey, 1*60, false)
	c6_ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Glimbright Shade (C6)",
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       c6_glimbrightAtkP,
	}
	c.Core.QueueAttack(
		c6_ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4),
		1, //TODO: c6 hitmark
		1,
	)
}
