package sandrone

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key     = "sandrone-c1"
	c2Key     = "sandrone-c2"
	c4Key     = "sandrone-c4"
	c4ICDKey  = "sandrone-c4-icd"
	c6Key     = "sandrone-c6"
	c6Hitmark = 5
)

// the first beam already has a stack so the 0.4 is never used
var c2Buff = []float64{0.4, 0.6, 0.8, 1.0}

func (c *char) c1DecoderGainMult() float64 {
	if c.Base.Cons < 1 {
		return 1.0
	}
	return 0.5
}

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c1Key, -1),
			Amount: func(ai info.AttackInfo) float64 {
				if c.currFagio != stateDecoding {
					return 0
				}

				switch ai.AttackTag {
				case attacks.AttackTagDirectStellarConduct:
				default:
					return 0
				}
				return 0.3
			},
		})
	}
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c2Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if c.getRadiance() == radianceNone {
				return nil
			}
			if !slices.Contains(atk.Info.AdditionalTags, attacks.AdditionalTagSandroneBeam) {
				return nil
			}
			m[attributes.CD] = c2Buff[c.c2stacks]
			return m
		},
	})
}

// stack gain before beam hits
func (c *char) c2OnBeam() {
	if c.Base.Cons < 2 {
		return
	}

	c.c2stacks = min(c.c2stacks+1, 3)
}

func (c *char) c2OnCAStart() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2stacks = 0
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		atk, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		var mult float64
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
			mult = 1.25
		case attacks.AttackTagDirectStellarSwirl:
			mult = 1.875
		default:
			return
		}

		if c.StatusIsActive(c4ICDKey) {
			return
		}

		c.AddStatus(c4ICDKey, 4*60, true)

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Prismatic Resonance Cannon (C4)",
			AttackTag:  atk.Info.AttackTag,
			Element:    attributes.Cryo,
			Mult:       mult,
		}

		ap := combat.NewBoxHitOnTarget(
			c.Core.Combat.Player(),
			info.Point{Y: -5},
			14,
			12,
		)

		c.Core.QueueAttack(ai, ap, 15, 15)
	}, c4Key)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.Core.Events.Subscribe(event.OnApplyAttack, func(args ...any) {
		atk := args[0].(*info.AttackEvent)

		// don't apply elevation to the reaction attack, since the subcomponent contributor attacks each got elevation applied already
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}
		atk.Info.Elevation += 0.2
	}, c6Key)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		if atk.Info.ActorIndex != c.Index() {
			return
		}

		atk.Info.Elevation += 0.2
	}, c6Key+"-stellar-atk")
}

func (c *char) c6OnCaStart() {
	if c.Base.Cons < 6 {
		return
	}
	c.c6stacks = 0
}

// stack gain before beam hits
func (c *char) c6OnBeam() {
	if c.Base.Cons < 6 {
		return
	}

	c.c6stacks += 1
	if c.c6stacks < 3 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex:     c.Index(),
		Abil:           "Condensed Cluster Beam (C6)",
		AttackTag:      attacks.AttackTagExtra,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagSandroneBeam},
		ICDTag:         attacks.ICDTagSandroneExtraAttackLaser,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Cryo,
		Durability:     25,
		Mult:           1,
	}

	switch c.getRadiance() {
	case radianceStellarConduct:
		ai.Abil += stellarConductText
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.ICDTag = attacks.ICDTagNone
		ai.Durability = 0
		ai.IgnoreDefPercent = 1
		ai.Mult = 0.8
	case radianceStellarSwirl:
		ai.Abil += stellarSwirlText
		ai.AttackTag = attacks.AttackTagDirectStellarSwirl
		ai.ICDTag = attacks.ICDTagNone
		ai.Durability = 0
		ai.IgnoreDefPercent = 1
		ai.Mult = 1.2
	}
	ap := combat.NewBoxHitOnTarget(
		c.Core.Combat.Player(),
		info.Point{Y: -7.5},
		3,
		15,
	)

	c.Core.QueueAttack(ai, ap, c6Hitmark, c6Hitmark)
}
