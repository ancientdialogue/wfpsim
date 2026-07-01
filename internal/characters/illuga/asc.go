package illuga

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const a1Key = "illuga-a1"

var (
	a4BonusGeo = []float64{0, 0.07, 0.14, 0.24, 0.24}
	a4BonusLCr = []float64{0, 0.48, 0.96, 1.6, 1.6}
)

func (c *char) a1Init() {
	if c.Base.Ascension < 1 {
		return
	}
	c.a1Buff = make([]float64, attributes.EndStatType)
	cr, cd := c.c6A1Buff()
	c.a1Buff[attributes.CR] = cr
	c.a1Buff[attributes.CD] = cd

	// Increases Elemental Mastery by 50
	c.a1BuffGleam = make([]float64, attributes.EndStatType)
	c.a1BuffGleam[attributes.EM] = c.c6A1BuffGleam()

	// workaround for giving lunarcrystallize the CR/CD
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if ae.Info.Element != attributes.Geo {
			return
		}

		if !c.Core.Player.ByIndex(ae.Info.ActorIndex).StatModIsActive(a1Key) {
			return
		}

		if c.Core.Flags.LogDebug {
			c.Core.Log.NewEvent("Illuga A1 added to Lunarcrystallize", glog.LogPreDamageMod, ae.Info.ActorIndex).
				Write("cr before", ae.Snapshot.Stats[attributes.CR]).
				Write("cr addition", c.a1Buff[attributes.CR]).
				Write("cd before", ae.Snapshot.Stats[attributes.CD]).
				Write("cd addition", c.a1Buff[attributes.CD])
		}
		ae.Snapshot.Stats[attributes.CR] += c.a1Buff[attributes.CR]
		ae.Snapshot.Stats[attributes.CD] += c.a1Buff[attributes.CD]
	}, a1Key+"-lcr")
}

func (c *char) a1OnSkillBurst() {
	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBaseWithHitlag(a1Key, 20*60),
			Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
				if atk.Info.Element != attributes.Geo {
					return nil
				}
				return c.a1Buff
			},
		})

		if c.Core.Player.GetMoonsignLevel() < 2 {
			continue
		}

		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(a1Key+"-gleam", 20*60),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return c.a1BuffGleam
			},
		})
	}
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Hydro:
		case attributes.Geo:
		default:
			continue
		}
		c.a4Count += 1
	}
}

func (c *char) a4BonusScaling(isLCr bool) float64 {
	if c.Base.Ascension < 4 {
		return 0
	}
	if isLCr {
		return a4BonusLCr[c.a4Count]
	}

	return a4BonusGeo[c.a4Count]
}
