package lohen

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	a4Key            = "lohen-a4"
	skillBonusKey    = "lohen-skill-bonus"
	skillBonusIcdKey = "lohen-skill-bonus-icd"
)

var cryoReactions = []event.Event{
	event.OnSwirlCryo,
	event.OnCrystallizeCryo,
	event.OnSuperconduct,
	event.OnFrozen,
	event.OnMelt,
	event.OnStellarSwirl,
	event.OnStellarConduct,
}

func (c *char) a1BonusWill(dmg float64) float64 {
	if c.Base.Ascension < 1 {
		return 0
	}
	if dmg < c.Stat(attributes.BaseATK)*30 {
		return 0
	}
	return 60
}

func (c *char) a4Init() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.ATKP] = 0.15
	statMod := character.StatMod{
		Base:         modifier.NewBaseWithHitlag(a4Key, 8*60),
		AffectedStat: attributes.ATKP,
		Amount: func() []float64 {
			return m
		},
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		ae := args[1].(*info.AttackEvent)
		index := ae.Info.ActorIndex
		if index == c.Index() {
			return
		}

		otherChar := c.Core.Player.Chars()[index]
		otherChar.AddStatMod(statMod)
		c.AddStatMod(statMod)
	}

	for _, evt := range cryoReactions {
		c.Core.Events.Subscribe(evt, hook, "lohen-a4-")
	}
}

func (c *char) TalentLvlSkill() int {
	if c.StatusIsActive(skillBonusKey) {
		return c.Character.TalentLvlSkill() + 1
	}
	return c.Character.TalentLvlSkill()
}

func (c *char) skillBonusOnSkillMasterstroke() {
	if c.StatusIsActive(skillBonusIcdKey) {
		return
	}
	c.AddStatus(skillBonusIcdKey, 18*60, true)
	dur := 9 * 60
	for _, char := range c.Core.Player.Chars() {
		if char.Index() == c.Index() {
			continue
		}
		charMaxLevel := max(char.TalentLvlAttack(), char.TalentLvlSkill(), char.TalentLvlBurst())
		if charMaxLevel >= c.Character.TalentLvlSkill() {
			dur += 6 * 60
			break
		}
	}
	c.AddStatus(skillBonusKey, dur, true)
}

func (c *char) hexereiInit() {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	hexBuff := make([]float64, attributes.EndStatType)
	hexBuff[attributes.DmgP] = 0.4
	c.hexereiAtkMod = character.AttackMod{
		Base: modifier.NewBaseWithHitlag("lohen-hexerei", 6*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag == attacks.AttackTagNormal || atk.Info.AttackTag == attacks.AttackTagExtra {
				return hexBuff
			}
			return nil
		},
	}
}

func (c *char) hexereiOnSkillBurst(will float64) {
	if !c.IsHexerei {
		return
	}

	if c.Core.Player.GetHexereiCount() < 2 {
		return
	}

	if will >= c.willToWinMax/2 {
		c.AddAttackMod(c.hexereiAtkMod)
	}
}
