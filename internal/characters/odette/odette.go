package odette

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func init() {
	core.RegisterCharFunc(keys.Odette, NewChar)
}

const (
	stellarConductText = " (Stellar-Conduct)"
	stellarSwirlText   = " (Stellar Swirl)"
)

type char struct {
	*tmpl.Character
	danceDoubleSrc int
	a1StacksSelf   int
	a1StacksOthers int
	a1Src          int
	c2Src          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.stellarInit()
	c.a1Init()
	c.c2Init()
	c.c4Init()
	c.c6Init()
	return nil
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 12
	}
	return c.Character.AnimationStartDelay(k)
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.StatusIsActive(danceDoubleKey) && !c.StatusIsActive(danceDoubleUpgradeKey) && !c.StatusIsActive(skillSpecialCDKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ResetActionCooldown(a action.Action) {
	if a == action.ActionSkill && c.StatusIsActive(danceDoubleKey) && !c.StatusIsActive(danceDoubleUpgradeKey) {
		c.DeleteStatus(skillSpecialCDKey)
		c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), "special skill cooldown forcefully reset").
			Write("type", a.String()).
			Write("charges_remain", c.AvailableCDCharge[a])
		return
	}
	c.Character.ResetActionCooldown(a)
}

func (c *char) ReduceActionCooldown(a action.Action, v int) {
	if a == action.ActionSkill && c.StatusIsActive(danceDoubleKey) && !c.StatusIsActive(danceDoubleUpgradeKey) {
		dur := max(c.StatusDuration(skillSpecialCDKey)-v, 0)
		c.AddStatus(skillSpecialCDKey, dur, false)
		c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index(), "special skill cooldown forcefully reduced").
			Write("type", a.String()).
			Write("charges_remain", c.AvailableCDCharge[a])
		return
	}
	c.Character.ReduceActionCooldown(a, v)
}

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *char) getRadiance() radianceState {
	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}
