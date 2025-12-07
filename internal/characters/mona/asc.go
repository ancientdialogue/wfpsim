package mona

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	phantasmalBubbleKey        = "phantasmal-bubble"
	phantasmalBubbleIcdKey     = "phantasmal-bubble-icd"
	hexereiOmenExtensionIcdKey = "hexerei-omen-extend-icd"
)

// After she has used Illusory Torrent for 2s, if there are any opponents nearby,
// Mona will automatically create a Phantom.
// A Phantom created in this manner lasts for 2s, and its explosion DMG is equal to 50% of Mirror Reflection of Doom.
//
// - checks for ascension level in dash.go to avoid queuing this up only to fail the ascension level check
func (c *char) a1() {
	// do nothing if not Mona
	if c.Core.Player.Active() != c.Index() {
		return
	}
	// do nothing if we aren't dashing anymore
	if c.Core.Player.CurrentState() != action.DashState {
		return
	}
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 15), nil)
	if enemies != nil {
		c.Core.Log.NewEvent("mona-a1 phantom added", glog.LogCharacterEvent, c.Index()).
			Write("expiry:", c.Core.F+120)
		// queue up phantom explosion
		phantomPos := c.Core.Combat.Player()
		c.Core.Tasks.Add(func() {
			aiExplode := info.AttackInfo{
				ActorIndex: c.Index(),
				Abil:       "Mirror Reflection of Doom (A1 Explode)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagNone,
				ICDGroup:   attacks.ICDGroupDefault,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Hydro,
				Durability: 25,
				Mult:       0.5 * skill[c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(aiExplode, combat.NewCircleHitOnTarget(phantomPos, nil, 5), 0, 0)
		}, 120)
	}
	// queue up next A1 check because Mona's still dashing
	// different Phantoms coexist and don't overwrite each other
	c.Core.Tasks.Add(c.a1, 120) // check again in 2s
}

// Increases Mona's Hydro DMG Bonus by a degree equivalent to 20% of her Energy Recharge rate.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}

	if c.a4Stats == nil {
		c.a4Stats = make([]float64, attributes.EndStatType)
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("mona-a4", -1),
			AffectedStat: attributes.HydroP,
			Extra:        true,
			Amount: func() ([]float64, bool) {
				return c.a4Stats, true
			},
		})
	}
	c.a4Stats[attributes.HydroP] = 0.2 * c.NonExtraStat(attributes.ER)
	c.QueueCharTask(c.a4, 60)
}

func (c *char) hexereiInit() {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("mona-hexerei-on-vaporize", -1),
			Amount: func(ai info.AttackInfo) (float64, bool) {
				if !ai.Amped {
					return 0, false
				}
				if ai.AmpType != info.ReactionTypeVaporize {
					return 0, false
				}
				if !c.StatusIsActive(phantasmalBubbleKey) {
					return 0, false
				}

				stacks := c.phantasmalBubbleStacks
				if stacks == 0 {
					return 0, false
				}
				c.phantasmalBubbleStacks = 0
				return 0.05 * float64(stacks), false
			},
		})
	}
}

func (c *char) makeHexereiCB() info.AttackCBFunc {
	if !c.IsHexerei {
		return nil
	}
	if c.Core.Player.GetHexereiCount() < 2 {
		return nil
	}

	return func(a info.AttackCB) {
		e, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		if !c.StatusIsActive(phantasmalBubbleIcdKey) {
			c.AddStatus(phantasmalBubbleIcdKey, 0.1*60, true)
			c.AddStatus(phantasmalBubbleKey, 8*60, true)
			c.phantasmalBubbleStacks = min(c.phantasmalBubbleStacks+1, 3)
		}

		omenExp := e.StatusExpiry(omenKey)
		if omenExp > c.Core.F && c.hexereiOmenExtension < 8*60 && !e.StatusIsActive(hexereiOmenExtensionIcdKey) {
			e.AddStatus(hexereiOmenExtensionIcdKey, 0.5*60, true)
			// calculate new duration
			newDur := omenExp - c.Core.F + 2*60
			e.AddStatus(omenKey, newDur, true)
			c.hexereiOmenExtension += 2 * 60
		}

		if e.StatusIsActive(bubbleKey) && !e.StatusIsActive(hexereiOmenExtensionIcdKey) {
			e.AddStatus(hexereiOmenExtensionIcdKey, 0.5*60, true)
			c.hexereiOmenExtension += 2 * 60
			c.omenStartingBonusDur = 2 * 60
		}
	}
}

func (c *char) hexereiOnBurst() {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	c.hexereiOmenExtension = 0
	c.omenStartingBonusDur = 0
}

func (c *char) getHexereiCount() int {
	count := 0
	for _, c := range c.Core.Player.Chars() {
		if c.IsHexerei {
			count += 1
		}
	}
	return count
}
