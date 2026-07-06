package venti

import (
	"fmt"

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
	hexereiBurstIcdKey = "hexerei-burst-icd"
	hexereiOnSwirlKey  = "hexerei-on-swirl"
)

var swirlEvents = []event.Event{
	event.OnSwirlPyro,
	event.OnSwirlCryo,
	event.OnStellarSwirl,
	event.OnSwirlElectro,
	event.OnSwirlHydro,
}

// A1 is not implemented and will likely never be implemented:
// Holding Skyward Sonnet creates an upcurrent that lasts for 20s.

// Regenerates 15 Energy for Venti after the effects of Wind's Grand Ode end.
// If an Elemental Absorption occurred, this also restores 15 Energy to all characters of that corresponding element in the party.
//
// - checks for ascension level in burst.go to avoid queuing this up only to fail the ascension level check
func (c *char) a4() {
	c.AddEnergy("venti-a4", 15)
	if c.qAbsorb == attributes.NoElement {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		if char.Base.Element == c.qAbsorb {
			char.AddEnergy("venti-a4", 15)
		}
	}
}

func (c *char) hexereiMakeBuff() func(args ...any) {
	return func(args ...any) {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return
		}

		if !c.StatusIsActive(burstKey) {
			return
		}
		c.AddStatus(hexereiOnSwirlKey, 4*60, true)

		c.Core.Player.ActiveChar().AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(hexereiOnSwirlKey+"-buff", 4*60),
			Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
				return c.hexereiBuff
			},
		})
	}
}

func (c *char) hexereiBurstBuff() float64 {
	if !c.IsHexerei {
		return 1.0
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return 1.0
	}

	if !c.StatusIsActive(hexereiOnSwirlKey) {
		return 1.0
	}

	return 1.35
}

func (c *char) hexereiInit() {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}
	c.hexereiBuff = make([]float64, attributes.EndStatType)
	c.hexereiBuff[attributes.DmgP] = 0.50
	for _, evt := range swirlEvents {
		c.Core.Events.Subscribe(evt, c.hexereiMakeBuff(), fmt.Sprintf("venti-hexerei-hook-%v", evt))
	}
}

func (c *char) hexereiNaBuff(mult float64) (info.AttackInfo, info.AttackPattern, info.AttackCBFunc) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  attacks.AttackTagNormal,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Physical,
		Mult:       mult,
		Durability: 25,
	}

	ap := combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		info.Point{Y: -0.5},
		0.1,
		1,
	)
	if !c.IsHexerei {
		return ai, ap, nil
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return ai, ap, nil
	}

	if !c.StatusIsActive(burstKey) {
		return ai, ap, nil
	}

	ai.Abil = fmt.Sprintf("Hurricane Arrow %v", c.NormalCounter)
	ai.Mult *= hurricaneBonus[c.TalentLvlAttack()]
	ai.ICDTag = attacks.ICDTagNormalAttack
	ai.Element = attributes.Anemo
	ai.IgnoreInfusion = true

	deltaPos := c.Core.Combat.Player().Pos().Sub(c.Core.Combat.PrimaryTarget().Pos())
	dist := deltaPos.Magnitude()

	// simulate piercing. Extends 15 units from target
	ap = combat.NewBoxHit(
		c.Core.Combat.Player(),
		c.Core.Combat.PrimaryTarget(),
		info.Point{Y: -dist},
		0.1,
		15,
	)

	return ai, ap, c.c2NaCB
}

func (c *char) hexereiNaCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.hexereiBurstExtCount >= 2 {
		return
	}

	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	dur := c.StatusDuration(burstKey)
	if dur == 0 {
		return
	}

	if c.StatusIsActive(hexereiBurstIcdKey) {
		return
	}

	c.AddStatus(hexereiBurstIcdKey, 0.1*60, true)
	c.AddStatus(burstKey, dur+60, false)
	c.qAbsorbBonusTicks += 3 // three extra ticks per extension?

	// extend CD?
	c.IncreaseActionCooldown(action.ActionBurst, 0.5*60)

	c.hexereiBurstExtCount += 1
}
