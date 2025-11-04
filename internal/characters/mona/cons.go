package mona

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c4key = "mona-c4"
	c6Key = "mona-c6"
)

// C1:
// When any of your own party members hits an opponent affected by an Omen, the effects of Hydro-related Elemental Reactions are enhanced for 8s:
// - Electro-Charged DMG increases by 15%.
// - Vaporize DMG increases by 15%.
// - Hydro Swirl DMG increases by 15%.
// - Frozen duration is extended by 15%.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}
	// TODO: "Frozen duration is extended by 15%." is bugged
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		// ignore if target doesn't have debuff
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.StatusIsActive(bubbleKey) && !t.StatusIsActive(omenKey) {
			return false
		}
		// add c1 to all party members, delay by 1, because:
		// "This bonus does not apply in the triggering attack nor from the resulting Hydro DMG dealt by Illusory Bubble in Stellaris Phantasm regardless if they were from resulting reactions."
		for _, x := range c.Core.Player.Chars() {
			char := x
			c.Core.Tasks.Add(func() {
				// TODO: "Vaporize DMG increases by 15%." should be getting snapshot, see https://library.keqingmains.com/evidence/characters/hydro/mona#mona-c1-snapshot-for-vape
				// requires ReactBonusMod refactor
				char.AddReactBonusMod(character.ReactBonusMod{
					Base: modifier.NewBase("mona-c1", 8*60),
					Amount: func(ai info.AttackInfo) (float64, bool) {
						bonus := 0.15
						if char.Index() != c.Core.Player.Active() {
							bonus *= 1.6
						}
						switch ai.AttackTag {
						// Hydro Swirl DMG increases by 15%.
						// Electro-Charged DMG increases by 15%.
						// Lunar-Charged DMG increases by 15%.
						case attacks.AttackTagSwirlHydro, attacks.AttackTagECDamage, attacks.AttackTagReactionLunarCharge, attacks.AttackTagDirectLunarCharged:
							return bonus, false
						}

						// Vaporize DMG increases by 15%.
						// the only way Hydro Swirl can vape is via an AoE Hydro Swirl which doesn't do damage anyways, so this is fine
						if ai.Amped && ai.AmpType == info.ReactionTypeVaporize {
							return bonus, false
						}

						return 0, false
					},
				})
			}, 1)
		}
		return false
	}, "mona-c1-check")
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2Buff = make([]float64, attributes.EndStatType)
	c.c2Buff[attributes.EM] = 80
}

func (c *char) c2OnBurst() {
	if c.Base.Cons < 2 {
		return
	}
	c.c2AfterBurst = true
}

// C2:
// When a Normal Attack hits, there is a 20% chance that it will be automatically followed by a Charged Attack.
// This effect can only occur once every 5s.
func (c *char) c2NaCB(a info.AttackCB) {
	trg := a.Target
	if c.Base.Cons < 2 {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.c2icd > c.Core.F {
		return
	}
	if !c.c2AfterBurst && c.Core.Rand.Float64() > .2 {
		return
	}
	c.c2AfterBurst = false
	c.c2icd = c.Core.F + 300 // every 5 seconds

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 10, 10, c.makeMagicCB(), c.c2CaCB, c.makeC6CAResetCB())
}

// C2:
// Additionally, when her Charged Attack hits an opponent, all nearby party members will have their Elemental Mastery increased by 80 for 12s.
func (c *char) c2CaCB(a info.AttackCB) {
	if c.Base.Cons < 2 {
		return
	}
	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("mona-c2", 480), // 8 s
			AffectedStat: attributes.EM,
			Amount: func() ([]float64, bool) {
				return c.c2Buff, true
			},
		})
	}
}

// C4:
// When any party member attacks an opponent affected by an Omen, their CRIT Rate is increased by 15%.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c4key, -1),
			Amount: func(_ *info.AttackEvent, t info.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				// exit if neither bubble nor omen are present
				if !x.StatusIsActive(bubbleKey) && !x.StatusIsActive(omenKey) {
					return nil, false
				}

				// Additionally, when any Magic party member attacks an opponent affected by an Omen, their CRIT DMG is increased by 15%.
				if char.IsMagic {
					m[attributes.CD] = 0.15
				} else {
					m[attributes.CD] = 0
				}

				return m, true
			},
		})
	}

	// workaround for giving lunarcharge the 15% CR
	c.Core.Events.Subscribe(event.OnLunarChargedReactionAttack, func(args ...any) bool {
		x, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return false
		}

		if !x.StatusIsActive(bubbleKey) && !x.StatusIsActive(omenKey) {
			return false
		}

		isMagic := c.Core.Player.ByIndex(ae.Info.ActorIndex).IsMagic

		if c.Core.Flags.LogDebug {
			evt := c.Core.Log.NewEvent("Mona C4 added to Lunarcharged", glog.LogPreDamageMod, ae.Info.ActorIndex).
				Write("before CR", ae.Snapshot.Stats[attributes.CR]).
				Write("additional CR", 0.15)
			if isMagic {
				evt.Write("before CDMG", ae.Snapshot.Stats[attributes.CD]).
					Write("additional CDMG", 0.15)
			}
		}

		ae.Snapshot.Stats[attributes.CR] += 0.15
		if isMagic {
			ae.Snapshot.Stats[attributes.CD] += 0.15
		}

		return false
	}, c4key+"-lunarcharged")
}

func (c *char) c6Check() bool {
	if c.Base.Cons < 6 {
		return false
	}

	monaDashing := c.Core.Player.Active() == c.Index() && c.Core.Player.CurrentState() == action.DashState

	// TODO: does this require her to be on field?
	nearbyOmen := false
	for _, e := range c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil) {
		if e.StatusIsActive(bubbleKey) || e.StatusIsActive(omenKey) {
			nearbyOmen = true
			break
		}
	}
	return monaDashing || nearbyOmen
}

func (c *char) c6() {
	if c.Base.Cons < 6 {
		return
	}
	// need to keep track of src in case of Mona Dash Dash, where the second dash starts between two c6 ticks
	// without a src check the second Dash would gain a stack before 1s is up and a second one at 1s

	// only start new task if a previous one isn't active
	if c.c6Src == -1 {
		c.c6Src = c.Core.F
		c.Core.Tasks.Add(c.c6Tick(c.Core.F), 60)
	}
}

// C6:
// Upon entering Illusory Torrent, Mona gains a 60% increase to the DMG of her next Charged Attack per second of movement.
// A maximum DMG Bonus of 180% can be achieved in this manner.
// The effect lasts for no more than 8s.
func (c *char) c6Tick(src int) func() {
	return func() {
		if c.c6Src != src {
			c.Core.Log.NewEvent(fmt.Sprintf("%v stack gain check ignored, src diff", c6Key), glog.LogCharacterEvent, c.Index()).
				Write("src", src).
				Write("new src", c.c6Src)
			return
		}

		if !c.c6Check() {
			c.c6Src = -1
			return
		}

		c.c6Stacks++
		if c.c6Stacks > 3 {
			c.c6Stacks = 3
		}
		c.Core.Log.NewEvent(fmt.Sprintf("%v stack gained", c6Key), glog.LogCharacterEvent, c.Index()).
			Write("c6Stacks", c.c6Stacks)

		m := make([]float64, attributes.EndStatType)
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c6Key, 8*60),
			Amount: func(atk *info.AttackEvent, t info.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagExtra {
					return nil, false
				}
				m[attributes.DmgP] = 0.60 * float64(c.c6Stacks)
				return m, true
			},
		})

		// reset C6 stacks in 8s if we didn't use a CA
		c.Core.Tasks.Add(c.c6TimerReset, 8*60+1)
		// queue up another stack and buff refresh in 1s
		c.Core.Tasks.Add(c.c6Tick(src), 60)
	}
}

func (c *char) makeC6CAResetCB() info.AttackCBFunc {
	if c.Base.Cons < 6 {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}

		if !c.StatusIsActive(c6Key) {
			return
		}
		c.DeleteStatus(c6Key)
		c.c6Stacks = 0
		c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via charge attack", c6Key), glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) c6TimerReset() {
	// handle C6 stack reset if CA not used before c6 buff expires
	if c.c6Stacks > 0 && !c.StatusIsActive(c6Key) {
		c.c6Stacks = 0
		c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via timer", c6Key), glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) bool {
		e, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}

		ae := args[1].(*info.AttackEvent)

		if ae.Info.ActorIndex != c.Index() {
			return false
		}

		if ae.Info.AttackTag != attacks.AttackTagExtra {
			return false
		}

		if !e.StatusIsActive(bubbleKey) && !e.StatusIsActive(omenKey) {
			return false
		}

		ae.Info.Mult *= 2.0
		return false
	}, "mona-c6-ca-omen")
}
