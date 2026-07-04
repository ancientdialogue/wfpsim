package qiqi

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	talismanKey        = "qiqi-talisman"
	talismanICDKey     = "qiqi-talisman-icd"
	stellarConductText = " (Stellar-Conduct)"
	radianceSwirlKey   = "radiance-stellar-swirl"
)

func init() {
	core.RegisterCharFunc(keys.Qiqi, NewChar)
}

type char struct {
	*tmpl.Character
	revelation        bool
	skillCD           int
	skillLastUsed     int
	skillHealSnapshot info.Snapshot // Required as both on hit procs and continuous healing need to use this
	c6Stacks          int
}

// TODO: Not implemented - C6 (revival mechanic, not suitable for sim)
// C4 - Enemy Atk reduction, not useful in this sim version
func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.skillCD = 30 * 60
	c.skillLastUsed = 0

	revelation, ok := p.Params["revelation"]
	if !ok {
		revelation = 1
	}
	c.revelation = revelation > 0

	w.Character = &c

	return nil
}

// Ensures the set of targets are initialized properly
func (c *char) Init() error {
	c.skillInit()
	c.revelationInit()
	c.a1()
	c.talismanHealHook()
	c.onNACAHitHook()
	c.c2Init()
	c.c6Init()
	return nil
}

// Helper function to calculate healing amount dynamically using current character stats, which has all mods applied
func (c *char) healDynamic(healScalePer, healScaleFlat []float64, talentLevel int) float64 {
	atk := c.TotalAtk()
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}

// Helper function to calculate healing amount from a snapshot instance
func (c *char) healSnapshot(d *info.Snapshot, healScalePer, healScaleFlat []float64, talentLevel int) float64 {
	atk := d.Stats.TotalATK()
	heal := healScaleFlat[talentLevel] + atk*healScalePer[talentLevel]
	return heal
}

func (c *char) AnimationStartDelay(k info.AnimationDelayKey) int {
	if k == info.AnimationXingqiuN0StartDelay {
		return 7
	}
	return c.Character.AnimationStartDelay(k)
}

type radianceState int

const (
	radianceNone radianceState = iota
	radianceStellarConduct
	radianceStellarSwirl
)

func (c *char) getRadiance() radianceState {
	if !c.revelation {
		return radianceNone
	}

	if c.StatusIsActive(reactable.PolestarFieldKey) {
		return radianceStellarConduct
	}

	if c.StatusIsActive(radianceSwirlKey) {
		return radianceStellarSwirl
	}

	return radianceNone
}
