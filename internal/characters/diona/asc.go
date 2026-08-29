package diona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

const (
	radianceSwirlKey   = "radiance-stellar-swirl"
	revelationSkillKey = "diona-revelation-skill"
	revelationKey      = "diona-revelation"
	revelationICDKey   = "diona-revelation-icd"
)

// Characters shielded by Icy Paws have their Movement SPD increased by 10% and their Stamina Consumption decreased by 10%.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	c.Core.Player.AddStamPercentMod("diona-a1", -1, func(_ action.Action) (float64, bool) {
		if c.Core.Player.Shields.Get(shield.DionaSkill) != nil {
			return -0.1, false
		}
		return 0, false
	})
}

// A4 is not implemented:
// TODO: Opponents who enter the AoE of Signature Mix have 10% decreased ATK for 15s.

func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)
	}, "diona-"+radianceSwirlKey)

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if !c.StatusIsActive(revelationSkillKey) {
			return
		}

		if c.StatusIsActive(revelationICDKey) {
			return
		}

		c.AddStatus(revelationICDKey, 3.5*60, true)

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Icy Paw (Revelation)",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       paw[c.TalentLvlSkill()],
		}

		for i := range 3 {
			c.Core.QueueAttack(
				ai,
				combat.NewCircleHit(
					c.Core.Combat.Player(),
					c.Core.Combat.PrimaryTarget(),
					nil,
					0.5,
				),
				0,
				5+i,
			)
		}
	}

	c.Core.Events.Subscribe(event.OnStellarSwirl, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnSwirlCryo, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnSuperconduct, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, revelationKey)
}

func (c *char) revelationOnSkill() {
	if !c.revelation {
		return
	}
	c.AddStatus(revelationSkillKey, 20*60, true)
}
