package columbina

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Columbina, NewChar)
}

type char struct {
	*tmpl.Character
	skillSrc            int
	a4stacks            int
	gravity             [2]float64
	gravityTask         bool
	gravityLastReaction info.ReactionType
	burstArea           info.AttackPattern
	a1Stacks            int
	a1Buff              []float64

	c2Buff []float64
	c4Buff []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = 3
	c.BurstCon = 3
	c.SkillCon = 5
	c.Moonsign = 1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillInit()
	c.a1Init()
	c.a4Init()
	c.moonsignInit()

	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 7
	}
	return c.Character.AnimationStartDelay(k)
}

// func (c *char) getMoonsignLevel() int {
// 	count := 0
// 	for _, c := range c.Core.Player.Chars() {
// 		count += c.Moonsign
// 	}
// 	return count
// }
