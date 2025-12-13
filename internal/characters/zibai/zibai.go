package zibai

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Zibai, NewChar)
}

type char struct {
	*tmpl.Character
	radiance   float64
	skillSrc   int
	skillsUsed int
	c6Elev     float64
}

const lunarCrystallizeAbil = " (Lunar-Crystallize)"

// TODO: need to clean up zhongli code still
func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = normalHitNum

	w.Character = &c

	c.Moonsign = 1

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	c.skillInit()
	c.a1Init()
	c.a4Init()
	c.moonsignInit()
	c.c1Init()
	c.c2Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 9
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if c.StatusIsActive(skillKey) && a == action.ActionSkill {
		if c.radiance < 70 {
			return false, action.InsufficientStamina
		}
		return true, action.NoFailure
	}

	return c.Character.ActionReady(a, p)
}

func (c *char) ResetNormalCounter() {
	c.c4ResetNormalCount()
}
