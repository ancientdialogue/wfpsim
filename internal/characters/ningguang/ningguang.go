package ningguang

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const skillParticleICDKey = "ningguang-skill-particle-ICD"

func init() {
	core.RegisterCharFunc(keys.Ningguang, NewChar)
}

type char struct {
	*tmpl.Character
	c2reset       int
	lastScreen    int
	skillSnapshot combat.Snapshot
}

func NewChar(s *core.Core, w *character.CharWrapper, _ character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	// Initialize at some very low value so these happen correctly at start of sim
	c.c2reset = -9999

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.Tags["jade"] > 0 {
			return 0
		}
		return 50
	}
	return c.Character.ActionStam(a, p)
}
