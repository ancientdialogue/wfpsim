package vodyanitsa

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	a1Key = "vodyanitsa-a1"
	a4Key = "vodyanitsa-a4"
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}

	// TODO: should be when the gadget is added
	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		if c.Core.Status.Duration(reactable.SswKey) > 0 {
			// we already had a vortex so no shred
			return
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
		for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
			e, ok := e.(*enemy.Enemy)
			if !ok {
				continue
			}
			e.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag(a1Key, 6*60),
				Ele:   attributes.Anemo,
				Value: -0.30,
			})
		}
	}, a1Key)

	// TODO: should be when the gadget is removed
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)

		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
		for _, e := range c.Core.Combat.EnemiesWithinArea(ap, nil) {
			e, ok := e.(*enemy.Enemy)
			if !ok {
				continue
			}
			e.AddResistMod(info.ResistMod{
				Base:  modifier.NewBaseWithHitlag(a1Key, 6*60),
				Ele:   attributes.Anemo,
				Value: -0.30,
			})
		}
	}, a1Key)
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		ae := args[1].(*info.AttackEvent)

		if !c.StatusIsActive(a4Key) {
			return
		}

		var amt float64
		switch ae.Info.AttackTag {
		case attacks.AttackTagElementalBurst, attacks.AttackTagElementalArt, attacks.AttackTagElementalArtHold, attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge:
			if ae.Info.Element != attributes.Cryo && ae.Info.Element != attributes.Hydro {
				return
			}

			if c.recentSSW() {
				return
			}
			amt = max(min((c.MaxHP()-40000)*140/1000, 3500), 0)
		case attacks.AttackTagDirectStellarSwirl:
			amt = max(min((c.MaxHP()-40000)*260/1000, 6500), 0)
		default:
			return
		}

		if ae.Info.ActorIndex == c.Core.Player.Active() {
			if c.leadVocal <= 0 {
				return
			}
			c.leadVocal--
		} else {
			if c.chorus <= 0 {
				return
			}
			c.chorus--
		}

		ae.Info.FlatDmg += amt
	}, a4Key)
}

func (c *char) a4OnSkill() {
	if c.Base.Ascension < 4 {
		return
	}

	c.leadVocal = 25
	c.chorus = 10
	c.AddStatus(a4Key, 30*60, true)
}
