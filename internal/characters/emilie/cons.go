package emilie

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1ICDKey = "emilie-c1-icd"
	c1ICDDur = 2.9 * 60
)

func (c *char) c1() {
	if c.Base.Cons < 1 {
		return
	}
	if c.Core.Combat.DamageMode {
		m := make([]float64, attributes.EndStatType)
		m[attributes.DmgP] = 0.2
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("emilie-c1", -1),
			Amount: func(a *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				_, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				if a.Info.Abil != "Lumidouce Case Lv1" && a.Info.Abil != "Lumidouce Case Lv2" && a.Info.Abil != "Cleardew Cologne" {
					return nil, false
				}
				return m, true
			},
		})
	}
	c.Core.Events.Subscribe(event.OnBurning, func(args ...interface{}) bool {
		_, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if c.StatusIsActive(c1ICDKey) {
			return false
		}
		c.AddStatus(c1ICDKey, c1ICDDur, false)
		if c.lumidouceLvl < 4 {
			c.lumidouceLvl++
		}
		if c.lumidouceLvl == 4 {
			c.lumidouceLvl = 2
			c.Core.Events.Emit(event.OnEmilieA1)
		}
		return false
	}, "emilie-c1-burning")

	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !t.IsBurning() {
			return false
		}
		if !ok {
			return false
		}
		if atk.Info.Element != attributes.Dendro {
			return false
		}
		if c.StatusIsActive(c1ICDKey) {
			return false
		}
		c.AddStatus(c1ICDKey, c1ICDDur, false)
		if c.lumidouceLvl < 4 {
			c.lumidouceLvl++
		}
		if c.lumidouceLvl == 4 {
			c.lumidouceLvl = 2
			c.Core.Events.Emit(event.OnEmilieA1)
		}
		return false
	}, "emilie-c1-on-damge")
}

func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.Abil != "Lumidouce Case Lv1" && atk.Info.Abil != "Lumidouce Case Lv2" && atk.Info.Abil != "Cleardew Cologne" {
			return false
		}

		t.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag("emilie-c2-shred", 10*60),
			Ele:   attributes.Dendro,
			Value: -0.3,
		})

		return false
	}, "emilie-c2")
}

func (c *char) c4Dur() int {
	if c.Base.Cons < 4 {
		return 0
	}

	return 120
}

func (c *char) c4Interval() int {
	if c.Base.Cons < 4 {
		return 0
	}

	return -0.3 * 60
}
