package cryo

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a4Key            = "travelercryo-a4"
	stellarBonusKey  = "travelercryo-stellar-bonus"
	radianceSwirlKey = "radiance-stellar-swirl"
)

func (c *Traveler) a1Conversion(ai *info.AttackInfo) {
	if c.Base.Ascension < 1 {
		return
	}

	if c.getRadiance() != radianceStellarConduct {
		return
	}

	if !c.StatusIsActive(skillKey) {
		return
	}

	ai.Element = attributes.Cryo
	ai.IgnoreInfusion = true
	ai.Mult += 0.8
}

func (c *Traveler) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(a4Key, -1),
		Extra:        true,
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			stats := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
			m[attributes.EM] = min(stats.TotalATK()*0.08, 160)
			return m
		},
	})
}

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *Traveler) getRadiance() radianceState {
	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}

func (c *Traveler) stellarInit() {
	reactable.EnableStellarConduct(c.Core)
	reactable.EnableStellarSwirl(c.Core)

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct:
		case attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.0035, 0.7)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("travelercryo adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		bonus := min(c.TotalAtk()/100.0*0.0035, 0.7)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("travelercryo adding stellar base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
	}, stellarBonusKey+"-reaction")

	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		c.AddStatus(radianceSwirlKey, 8*60, false)
	}, stellarBonusKey)
}
