package columbina

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
	lunarchargeBonusKey = "columbina-lc-bonus"
	a1Key               = "columbina-a1"
)

func (c *char) lunarchargeInit() {
	c.Core.Flags.Custom[reactable.LunarChargeEnableKey] = 1

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		atk := args[1].(*info.AttackEvent)

		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectLunarCharged:
		case attacks.AttackTagReactionLunarCharge:
		default:
			return false
		}

		bonus := min(c.MaxHP()/1000.0*0.002, 0.07)

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("columbina adding lunarcharged base damage", glog.LogCharacterEvent, c.Index()).Write("bonus", bonus)
		}

		atk.Info.BaseDmgBonus += bonus
		return false
	}, lunarchargeBonusKey)
}

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Buff = make([]float64, attributes.EndStatType)
}

func (c *char) a1OGravityTick() {
	if c.Base.Ascension < 1 {
		return
	}
	if !c.StatModIsActive(a1Key) {
		c.a1Stacks = 1
	} else {
		c.a1Stacks = min(c.a1Stacks+1, 3)
	}

	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a1Key, 10*60),
		AffectedStat: attributes.CR,
		Amount: func() ([]float64, bool) {
			c.a1Buff[attributes.CR] = 0.05 * float64(c.a1Stacks)
			return c.a1Buff, true
		},
	})
}

func (c *char) a4Init() {
	a4Hook := func(args ...any) bool {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return false
		}

		if c.StatusIsActive(burstBuffKey) && c.Core.Combat.Player().IsWithinArea(c.burstArea) {
			c.Core.Flags.Custom[reactable.LcIcdOverrideKey] = 1.5 * 60
			c.Core.Flags.Custom[reactable.LcrExtraHitOverride] = 0.33
			return false
		}

		// player is outside of lunar domain, reset buffs
		delete(c.Core.Flags.Custom, reactable.LcIcdOverrideKey)
		delete(c.Core.Flags.Custom, reactable.LcrExtraHitOverride)
		return false
	}

	c.Core.Events.Subscribe(event.OnLunarCharged, a4Hook, "columbina-gravity-lc")
}
