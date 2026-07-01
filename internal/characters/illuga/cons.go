package illuga

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
	c1Key    = "illuga-c1"
	c1IcdKey = "illuga-c1-icd"
	c2Key    = "illuga-c2"
	c2IcdKey = "illuga-c2-icd"
	c4Key    = "illuga-c4"
	c4IcdKey = "illuga-c4-icd"
	c6Key    = "illuga-c6"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if c.StatusIsActive(c1IcdKey) {
			return
		}
		c.AddStatus(c1IcdKey, 15*60, true)
		c.AddEnergy(c1Key, 12)
	}

	c.Core.Events.Subscribe(event.OnCrystallizeCryo, hook, c1Key+"-cryo")
	c.Core.Events.Subscribe(event.OnCrystallizeHydro, hook, c1Key+"-hydro")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, hook, c1Key+"-electro")
	c.Core.Events.Subscribe(event.OnCrystallizePyro, hook, c1Key+"-pyro")
	c.Core.Events.Subscribe(event.OnLunarCrystallize, hook, c1Key+"-lcr")
}

func (c *char) c2OnStackConsume() {
	if c.Base.Cons < 2 {
		return
	}

	stacksUsed := initalStacks + c.stacksGained - c.burstStacks
	if stacksUsed%7 > 0 {
		return
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Illuga C2",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Geo,
		Durability: 25,
		UseDef:     true,
		Mult:       2.0,
		FlatDmg:    4.0 * c.Stat(attributes.EM),
	}

	c.Core.QueueAttack(
		ai,
		combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()),
		3,
		3,
	)
}

func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.DEF] = 200

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c4Key, -1),
			AffectedStat: attributes.DEF,
			Amount: func() []float64 {
				if c.Core.Player.Active() != char.Index() {
					return nil
				}

				if !c.StatusIsActive(burstKey) {
					return nil
				}
				return m
			},
		})
	}
}

func (c *char) c6A1Buff() (float64, float64) {
	if c.Base.Cons < 6 {
		return 0.05, 0.1
	}
	return 0.1, 0.3
}

func (c *char) c6A1BuffGleam() float64 {
	if c.Base.Cons < 6 {
		return 50
	}
	return 80
}
