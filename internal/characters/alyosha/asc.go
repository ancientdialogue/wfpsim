package alyosha

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const a4Key = "alyosha-a4"

func (c *char) a1OnBurstTick() {
	if c.Base.Ascension < 1 {
		return
	}

	c.Core.Player.Heal(info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Alyosha (C1)",
		Src:     1.2 * c.TotalAtk(),
		Bonus:   c.Stat(attributes.Heal),
	})
}

func (c *char) a4Init() {
	if c.Base.Ascension < 4 {
		return
	}

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(a4Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagElementalArt:
			case attacks.AttackTagElementalArtHold:
			case attacks.AttackTagElementalBurst:
			default:
				return nil
			}
			m[attributes.DmgP] = min(c.Stat(attributes.ER)/0.01*0.0035, 0.7)
			return m
		},
	})
}

func (c *char) stellarInit() {
	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("alyosha-ssc", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.isStellarRadiance() {
					return 0
				}
				if !c.StatusIsActive(skillBuffKey) {
					return 0
				}
				if ai.AttackTag != attacks.AttackTagDirectStellarConduct {
					return 0
				}
				return 0.2 * float64(c.skillStacks)
			},
		})
	}
}

func (c *char) isStellarRadiance() bool {
	return c.StatusIsActive(reactable.PolestarFieldKey)
}
