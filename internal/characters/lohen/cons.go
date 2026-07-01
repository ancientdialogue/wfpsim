package lohen

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c2Key    = "lohen-c2"
	c2IcdKey = "lohen-c2-icd"
	c4Key    = "lohen-c4"
	c6Key    = "lohen-c6"
	c6IcdKey = "lohen-c6-icd"
)

func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.willToWinMax = 300
}

func (c *char) c1WillToWinMult() float64 {
	if c.Base.Cons < 1 {
		return 1
	}
	return 5
}

func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 200
	c.c2StatMod = character.StatMod{
		Base: modifier.NewBaseWithHitlag("lohen-c2-em", 8*60),
		Amount: func() []float64 {
			return m
		},
	}
}

// When in Masterstroke mode, after unleashing the special Elemental Skill Etched Into Bone and Soul or the Elemental Burst Manifest Judgment, Lohen gains "Evilsbane Blade" for 4s: The next time Lohen hits an opponent with a Normal or Charged Attack while in Masterstroke mode, he will follow up with an additional strike that deals AoE Cryo DMG equal to 300% of his ATK, and increase the Elemental Mastery of other nearby party members by 125 for 8s. Evilsbane Blade can be triggered once every 4s.
func (c *char) c2OnSkillBurst() {
	if c.Base.Cons < 2 {
		return
	}
	if c.StatusIsActive(c2IcdKey) {
		return
	}

	c.AddStatus(c2IcdKey, 4*60, true)
	c.AddStatus(c2Key, 4*60, true)
}

func (c *char) c2MakeCB() info.AttackCBFunc {
	if c.Base.Cons < 2 {
		return nil
	}

	return func(a info.AttackCB) {
		if !c.StatusIsActive(c2Key) {
			return
		}
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Lohen C2",
			AttackTag:  attacks.AttackTagNormal,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Cryo,
			Durability: 25,
			Mult:       5,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(a.Target, nil, 4), 5, 5)
		for _, char := range c.Core.Player.Chars() {
			if char.Index() == c.Index() {
				continue
			}
			char.AddStatMod(c.c2StatMod)
		}
		c.DeleteStatus(c2Key)
	}
}

func (c *char) c4OnBurst() {
	if c.Base.Cons < 4 {
		return
	}
	if !c.StatusIsActive(skillKey) {
		return
	}
	c.gainWillToWin(c.willToWinMax, "lohen-c4")
}

func (c *char) c4OnBurstEnergyConsume() {
	if c.Base.Cons < 4 {
		return
	}
	if !c.StatusIsActive(c4Key) {
		return
	}
	c.AddEnergy("lohen-c4", 15)
}

func (c *char) c4OnSkillMasterstroke() {
	if c.Base.Cons < 4 {
		return
	}
	if c.Energy < c.EnergyMax {
		c.AddEnergy("lohen-c4", 15)
		return
	}
	c.AddStatus(c4Key, 15*60, true)
}

func (c *char) c6ExtendOnSkill() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.StatusIsActive(c6Key) {
		return
	}
	c.ExtendStatus(skillKey, 1.25*60)
	c.maxSkillsUsed += 1
}

func (c *char) c6OnSkillBurst(amt float64, snap *info.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}
	if !c.StatusIsActive(skillKey) {
		return
	}

	if c.StatusIsActive(c6Key) {
		c.c6AddCD(snap)
		c.DeleteStatus(c6Key)
		return
	}

	if c.StatusIsActive(c6IcdKey) {
		return
	}

	// reset will to win to the amount consumed
	c.willToWin = amt
	c.joy = joyMax
	c.AddStatus(c6Key, 18*60, true)
	c.AddStatus(c6IcdKey, 7*60, true)
	c.c6AddCD(snap)
}

func (c *char) c6AddCD(snap *info.Snapshot) {
	if c.Base.Cons < 6 {
		return
	}
	old := snap.Stats[attributes.CD]
	snap.Stats[attributes.CD] += 1.75
	c.Core.Log.NewEvent("c6 adding crit damage", glog.LogCharacterEvent, c.Index()).
		Write("old", old).
		Write("new", snap.Stats[attributes.CD])
}
