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
	magicBurstIcdKey = "magic-burst-icd"
	magicOnSwirlKey  = "magic-on-swirl"
)

var swirlEvents = []event.Event{
	event.OnSwirlPyro,
	event.OnSwirlCryo,
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

func (c *char) magicMakeBuff() func(args ...any) bool {
	return func(args ...any) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Core.Player.Active() {
			return false
		}

		if !c.StatusIsActive(burstKey) {
			return false
		}
		c.AddStatus(magicOnSwirlKey, 4*60, true)

		c.Core.Player.ActiveChar().AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(magicOnSwirlKey+"-buff", 4*60),
			Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
				return c.magicBuff, true
			},
		})

		return false
	}
}

func (c *char) magicBurstBuff() float64 {
	if !c.IsMagic {
		return 1.0
	}

	if c.getMagicCount() < 2 {
		return 1.0
	}

	if !c.StatusIsActive(magicOnSwirlKey) {
		return 1.0
	}

	return 1.35
}

func (c *char) magicInit() {
	if !c.IsMagic {
		return
	}

	if c.getMagicCount() < 2 {
		return
	}
	c.magicBuff = make([]float64, attributes.EndStatType)
	c.magicBuff[attributes.DmgP] = 0.50
	for _, evt := range swirlEvents {
		c.Core.Events.Subscribe(evt, c.magicMakeBuff(), fmt.Sprintf("venti-magic-hook-%v", evt))
	}
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

func (c *char) magicNaBuff(mult float64) (info.AttackInfo, info.AttackPattern, info.AttackCBFunc) {
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
	if !c.IsMagic {
		return ai, ap, nil
	}

	if c.getMagicCount() < 2 {
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

func (c *char) magicNaCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.magicBurstExtCount >= 2 {
		return
	}

	if !c.IsMagic {
		return
	}

	if c.getMagicCount() < 2 {
		return
	}

	dur := c.StatusDuration(burstKey)
	if dur == 0 {
		return
	}

	if c.StatusIsActive(magicBurstIcdKey) {
		return
	}

	c.AddStatus(magicBurstIcdKey, 0.1*60, true)
	c.AddStatus(burstKey, dur+60, false)
	c.qAbsorbBonusTicks += 3 // three extra ticks per extension?

	// extend CD?
	c.IncreaseActionCooldown(action.ActionBurst, 0.5*60)

	c.magicBurstExtCount += 1
}
