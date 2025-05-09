package skirk

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Skirk, NewChar)
}

type char struct {
	*tmpl.Character
	serpentsSubtlety float64
	skillSrc         int
	burstCount       int
	burstVoids       int
	voidRiftCount    int
	a4Stacks         []int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 0
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.serpentsSubtlety = maxSerpentsSubtlety
	c.onExitField()
	c.BurstInit()
	c.a1Init()
	c.a4Init()
	c.talentPassiveInit()
	return nil
}

func (c *char) talentPassiveInit() {
	allCryoHydro := true
	hasCryo := false
	hasHydro := false

	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Hydro:
			hasHydro = true
		case attributes.Cryo:
			hasCryo = true
		default:
			allCryoHydro = false
		}
	}
	if !allCryoHydro {
		return
	}
	if !hasCryo {
		return
	}
	if !hasHydro {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.SetTag(keys.SkirkPassive, 1)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if k == model.AnimationXingqiuN0StartDelay {
		return 6
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	switch a {
	case action.ActionBurst:
		if !c.StatusIsActive(skillKey) && c.serpentsSubtlety < 50 {
			return false, action.InsufficientEnergy
		}
		return c.Character.ActionReady(a, p)
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) NextQueueItemIsValid(k keys.Char, a action.Action, p map[string]int) error {
	// can use charge without attack beforehand unlike most of the other sword users
	if a == action.ActionCharge {
		return nil
	}
	return c.Character.NextQueueItemIsValid(k, a, p)
}

func (c *char) AddSerpentsSubtlety(src string, e float64) {
	pre := c.serpentsSubtlety
	c.serpentsSubtlety += e
	if c.serpentsSubtlety > maxSerpentsSubtlety {
		c.serpentsSubtlety = maxSerpentsSubtlety
	}
	if c.serpentsSubtlety < 0 {
		c.serpentsSubtlety = 0
	}

	c.Core.Log.NewEvent(fmt.Sprintf("+%.1f serpent's subtlety, next: %.1f", e, c.serpentsSubtlety), glog.LogEnergyEvent, c.Index).
		Write("added", e).
		Write("pre_recovery", pre).
		Write("post_recovery", c.serpentsSubtlety).
		Write("source", src)
}

func (c *char) ReduceSerpentsSubtlety(src string, e float64) {
	pre := c.serpentsSubtlety
	c.serpentsSubtlety -= e
	if c.serpentsSubtlety > maxSerpentsSubtlety {
		c.serpentsSubtlety = maxSerpentsSubtlety
	}
	if c.serpentsSubtlety < 0 {
		c.serpentsSubtlety = 0
	}

	c.Core.Log.NewEvent(fmt.Sprintf("-%.1f serpent's subtlety, next: %.1f", e, c.serpentsSubtlety), glog.LogEnergyEvent, c.Index).
		Write("reduced", e).
		Write("pre", pre).
		Write("post", c.serpentsSubtlety).
		Write("source", src)
}

func (c *char) ConsumeSerpentsSubtlety(delay int, src string) {
	if delay == 0 {
		c.Core.Log.NewEvent("draining serpent's subtlety", glog.LogEnergyEvent, c.Index).
			Write("pre_drain", c.serpentsSubtlety).
			Write("post_drain", 0).
			Write("source", src)
		c.serpentsSubtlety = 0
		return
	}
	c.Core.Tasks.Add(func() {
		c.Core.Log.NewEvent("draining serpent's subtlety", glog.LogEnergyEvent, c.Index).
			Write("pre_drain", c.serpentsSubtlety).
			Write("post_drain", 0).
			Write("source", src)
		c.serpentsSubtlety = 0
	}, delay)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatusIsActive(skillKey) {
			c.exitSkillState(c.skillSrc)
		}
		return false
	}, "skirk-exit")
}
