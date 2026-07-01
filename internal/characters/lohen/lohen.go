package lohen

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Lohen, NewChar)
}

type char struct {
	*tmpl.Character
	joy           int
	willToWin     float64
	willToWinMax  float64
	skillsUsed    int
	maxSkillsUsed int
	hexereiAtkMod character.AttackMod
	c2StatMod     character.StatMod
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.SkillCon = 3
	c.BurstCon = 5

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum

	c.willToWinMax = 100

	w.Character = &c

	hexerei, ok := p.Params["hexerei"]
	if !ok {
		hexerei = 1
	}
	c.IsHexerei = hexerei > 0

	return nil
}

func (c *char) Init() error {
	c.hexereiInit()
	c.skillInit()
	c.onExitField()
	c.a4Init()
	c.c1Init()
	c.c2Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 11
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(skillKey) {
		if c.joy >= joyMax {
			return true, action.NoFailure
		}
		return false, action.InsufficientStamina
	}

	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && c.StatusIsActive(skillKey) {
		return 10
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "joy":
		if c.StatusIsActive(skillKey) {
			return c.joy, nil
		}
		return 0, nil
	case "will-to-win":
		if c.StatusIsActive(skillKey) {
			return c.willToWin, nil
		}
		return 0, nil
	default:
		return c.Character.Condition(fields)
	}
}
