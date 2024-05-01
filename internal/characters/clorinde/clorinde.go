package clorinde

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Clorinde, NewChar)
}

type char struct {
	*tmpl.Character

	a1stacks *stackTracker
	a4stacks *stackTracker
	a4bonus  []float64

	// track bol manually skip template
	hpDebt float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = base.SkillDetails.BurstEnergyCost
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.a1()
	c.a4Init()
	if c.Base.Cons >= 1 {
		c.c1()
	}
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if a1 window is active is on-field
	if a == action.ActionSkill && c.StatusIsActive(skillStateKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

// TODO: pew pew driver
func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 0
	case model.AnimationYelanN0StartDelay:
		return 0
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
