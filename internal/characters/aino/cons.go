package aino

import (
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
	c1Key    = "aino-c1"
	c2Key    = "aino-c2"
	c2IcdKey = "aino-c2-icd"
	c4Key    = "aino-c4"
	c4IcdKey = "aino-c4-icd"
	c6Key    = "aino-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.c1Buff = make([]float64, attributes.EndStatType)
	c.c1Buff[attributes.EM] = 80

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c1Key+"-buff", -1),
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.Active() != char.Index() {
					return nil, false
				}
				if !c.StatusIsActive(c1Key) {
					return nil, false
				}
				return c.c1Buff, true
			},
		})
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase(c1Key+"-buff", -1),
		AffectedStat: attributes.EM,
		Amount: func() ([]float64, bool) {
			if !c.StatusIsActive(c1Key) {
				return nil, false
			}
			return c.c1Buff, true
		},
	})
}

func (c *char) c1OnSkillBurst() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddStatus(c1Key, 15*60, true)
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		atk := args[1].(*info.AttackEvent)
		if c.Core.Player.Active() == c.Index() {
			return false
		}
		if atk.Info.ActorIndex == c.Index() {
			return false
		}
		if !c.StatusIsActive(burstKey) {
			return false
		}
		if c.StatusIsActive(c2IcdKey) {
			return false
		}
		c.AddStatus(c2IcdKey, 5*60, true)

		em := c.Stat(attributes.EM)
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Cool Your Jets Ducky (C2)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 25,
			Mult:       0.25,
			FlatDmg:    em + c.a4Dmg(),
		}

		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(e, nil, 2.5), 0, 10)

		return false
	}, c2Key)
}

func (c *char) c4CB(a info.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(c4IcdKey) {
		return
	}
	c.AddStatus(c4IcdKey, 10*60, true)

	c.AddEnergy(c4Key, 10)
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	c.c1Buff = make([]float64, attributes.EndStatType)
	c.c1Buff[attributes.EM] = 80

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c6Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if !c.StatusIsActive(c6Key) {
					return 0, false
				}
				buff := 0.15
				if c.getMoonsignLevel() >= 2 {
					buff += 0.2
				}

				switch ai.AttackTag {
				case attacks.AttackTagECDamage:
				case attacks.AttackTagBloom:
				case attacks.AttackTagBountifulCore: // TODO: Does aino buff bountiful bloom damage?
				case attacks.AttackTagDirectLunarCharged:
				case attacks.AttackTagReactionLunarCharge:
				// TODO: Lunar-Bloom
				default:
					return 0, false
				}
				return buff, false
			},
		})
	}
}

func (c *char) c6OnBurst() {
	if c.Base.Cons < 6 {
		return
	}

	c.AddStatus(c6Key, 15*60, true)
}
