package durin

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames        []int
	skillRecastFrames  []int
	skillBlackHitmarks []int = []int{20, 20 + 5, 20 + 5 + 5}
)

const (
	particleICDKey     = "durin-particle-icd"
	skillWindowKey     = "durin-essential-transformation"
	skillWindowDur     = 3 * 60
	skillCDStarts      = 19
	skillWhiteHitmarks = 21

	whiteKey     = "confirmation-of-purity"
	blackKey     = "denial-of-darkness"
	energyIcdKey = "durin-skill-energy-icd"
)

func init() {
	// Tap E
	skillFrames = frames.InitAbilSlice(180)
	skillFrames[action.ActionAttack] = 42
	skillFrames[action.ActionSkill] = 30

	// Recast White
	skillRecastFrames = frames.InitAbilSlice(88)
	skillRecastFrames[action.ActionLowPlunge] = 52
	skillRecastFrames[action.ActionSkill] = 44
	skillRecastFrames[action.ActionWalk] = 86
	skillRecastFrames[action.ActionSwap] = 87
}

// Dashes nimbly forward with silken steps. Once this dash ends, Chiori will
// summon the automaton doll "Tamoto" beside her and sweep her blade upward,
// dealing AoE Geo DMG to nearby opponents based on her ATK and DEF. Holding the
// Skill will cause it to behave differently.
//
// Hold Enter Aiming Mode to adjust the dash direction.
//
// Tamoto
// - Will slash at nearby opponents at intervals, dealing AoE Geo DMG based on
// Chiori's ATK and DEF.
// - While active, when Geo Construct(s) are created nearby, an additional Tamoto
// will be summoned next to Chiori. Only 1 additional Tamoto can be summoned in
// this manner, and its duration is independently counted.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	// if this is second press, swap and activate a1
	if c.StatusIsActive(skillWindowKey) {
		return c.skillRecastWhite()
	}

	c.AddStatus(skillWindowKey, skillWindowDur, true)
	c.SetCDWithDelay(action.ActionSkill, 12*60, skillCDStarts)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillRecastWhite() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Break: Lustrous Light",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       skillWhite[c.TalentLvlSkill()],
	}
	c.QueueCharTask(func() {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 0, 0, c.particleCB)
	}, skillWhiteHitmarks)
	c.DeleteStatus(skillWindowKey)
	c.DeleteStatus(blackKey)

	c.AddStatus(whiteKey, 30*60, true)

	if !c.StatusIsActive(energyIcdKey) {
		c.AddEnergy("durin-skill", skillEnergy[c.TalentLvlSkill()])
		c.AddStatus(energyIcdKey, 6*60, true)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) skillRecastBlack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Break: Dark Decay",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
	}
	for i, mult := range skillBlack {
		ai.Mult = mult[c.TalentLvlSkill()]
		c.QueueCharTask(func() {
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3), 0, 0, c.particleCB)
		}, skillBlackHitmarks[i])
	}
	c.DeleteStatus(skillWindowKey)
	c.DeleteStatus(whiteKey)

	c.AddStatus(blackKey, 30*60, true)
	if !c.StatusIsActive(energyIcdKey) {
		c.AddEnergy("durin-skill", skillEnergy[c.TalentLvlSkill()])
		c.AddStatus(energyIcdKey, 6*60, true)
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillRecastFrames),
		AnimationLength: skillRecastFrames[action.InvalidAction],
		CanQueueAfter:   skillRecastFrames[action.ActionSkill],
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 1*60, true)

	count := 3.0
	// if c.Core.Rand.Float64() < 0.33 {
	// 	count = 2
	// }
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Pyro, c.ParticleDelay)
}
