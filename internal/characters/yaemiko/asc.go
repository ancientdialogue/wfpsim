package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
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
	stellarConductText = " (Steller-Conduct)"
	revelationKey      = "yae-revelation"
	revelationICDKey   = "yae-revelation-icd"
)

// When casting Great Secret Art: Tenko Kenshin, each Sesshou Sakura destroyed
// resets the cooldown for 1 charge of Yakan Evocation: Sesshou Sakura.
func (c *char) a1OnBurstTotem() {
	if c.Base.Ascension < 1 {
		return
	}
	c.ResetActionCooldown(action.ActionSkill)
}

// When casting Great Secret Art: Tenko Kenshin, each Sesshou Sakura destroyed
// resets the cooldown for 1 charge of Yakan Evocation: Sesshou Sakura.
func (c *char) a1OnSkillPopKitsune() {
	if c.Base.Ascension < 1 {
		return
	}

	if !c.revelation {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Yae A4",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 0,
		Mult:       0.4,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)

	if c.isRadianceSSC() {
		ai.AttackTag = attacks.AttackTagDirectStellarConduct
		ai.Abil += stellarConductText
		ai.IgnoreDefPercent = 1
		ai.Mult = 0.5
	}

	c.Core.QueueAttack(ai, ap, 10, 10)
}

// Every point of Elemental Mastery Yae Miko possesses will increase Sesshou Sakura DMG by 0.15%.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("yaemiko-a4", -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			// only trigger on elemental art damage
			if atk.Info.AttackTag != attacks.AttackTagElementalArt {
				return nil
			}
			m[attributes.DmgP] = c.Stat(attributes.EM) * 0.0015
			return m
		},
	})
}

func (c *char) revelationInit() {
	if !c.revelation {
		return
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if c.revelationEnhanceTick {
			return
		}

		if c.StatusIsActive(revelationICDKey) {
			return
		}

		c.AddStatus(revelationICDKey, 2.5*60, true)

		// TODO: this only lasts 8s...
		c.revelationEnhanceTick = true
	}

	c.Core.Events.Subscribe(event.OnSuperconduct, hook, revelationKey)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, revelationKey)

	c.skillDur += 10 * 60
}

func (c *char) revelationEnhanceDMG() float64 {
	if !c.revelation {
		return 0
	}

	if !c.revelationEnhanceTick {
		return 0
	}

	c.revelationEnhanceTick = false

	dmg := 0.8 * c.TotalAtk()
	if !c.isRadianceSSC() {
		return dmg
	}

	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "Sesshou Sakura" + stellarConductText,
		AttackTag:        attacks.AttackTagDirectStellarConduct,
		ICDTag:           attacks.ICDTagNone,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Electro,
		Mult:             2,
		IgnoreDefPercent: 1,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 0.5)
	c.Core.QueueAttack(ai, ap, 10, 10)

	return dmg
}
