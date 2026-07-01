package zibai

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames        []int
	skillSkillFrames   []int
	skillStrideHitmark = []int{12, 24}
)

const (
	skillHitmark      = 19
	particleICDKey    = "zibai-particle-icd"
	skillKey          = "zibai-skill"
	radianceNAICDKey  = "zibai-radiance-na-icd"
	radianceLCrICDKey = "zibai-radiance-lcr-icd"
	maxRadiance       = 100

	skillAbil1 = "Spirit Steed's Stride"
	skillAbil2 = skillAbil1 + lunarCrystallizeAbil
)

func init() {
	skillFrames = frames.InitAbilSlice(30)

	skillSkillFrames = frames.InitAbilSlice(42)
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		// do nothing if previous char wasn't zibai
		prev := args[0].(int)
		if prev != c.Index() {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}

		c.DeleteStatus(skillKey)
	}, "zibai-exit")
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(skillKey) {
		return c.skillSkill()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Lunar Phase Shift (0 dmg)",
		AttackTag:  attacks.AttackTagNone,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Physical,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3), skillHitmark, skillHitmark)

	c.AddStatus(skillKey, 15*60+skillHitmark, true)

	c.SetCDWithDelay(action.ActionSkill, 18*60, skillHitmark)

	c.skillsUsed = 0

	c.skillSrc = c.Core.F

	c.QueueCharTask(func() {
		c.radianceTicker(c.skillSrc)()
		c.a1OnSkill()
		c.c1OnSkill()
	}, skillHitmark)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillSkill() (action.Info, error) {
	for i := range skillStride {
		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       skillAbil1,
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Geo,
			UseDef:     true,
			Durability: 25,
			Mult:       skillStride[i][c.TalentLvlSkill()],
		}

		if i == 1 {
			ai.Abil = skillAbil2
			ai.AttackTag = attacks.AttackTagDirectLunarCrystallize
			ai.Durability = 0
			ai.IgnoreDefPercent = 1
			ai.FlatDmg += c.a1StrideBonusDmg()
		}

		ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 5, 60)
		c.Core.QueueAttack(ai, ap, skillStrideHitmark[i], skillStrideHitmark[i], c.particleCB, c.c4SkillCB)
	}
	c.DeleteStatus(radianceLCrICDKey)
	c.skillsUsed += 1
	c.c6ConsumeRadiance()
	if c.skillsUsed >= c.c1MaxSkillsPerSkill() {
		c.DeleteStatus(skillKey)
	}

	return action.Info{
		Frames:          func(next action.Action) int { return skillSkillFrames[next] },
		AnimationLength: skillSkillFrames[action.InvalidAction],
		CanQueueAfter:   skillSkillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(particleICDKey) {
		return
	}

	if c.Core.Rand.Float64() > 0.67 {
		return
	}

	c.AddStatus(particleICDKey, 2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Geo, c.ParticleDelay)
}

func (c *char) radianceTicker(src int) func() {
	return func() {
		if c.skillSrc != src {
			return
		}

		if !c.StatusIsActive(skillKey) {
			return
		}

		c.addRadiance(1)
		c.Core.Tasks.Add(c.radianceTicker(src), 6)
	}
}

func (c *char) radianceCB(ac info.AttackCB) {
	if ac.Target.Type() != info.TargettableEnemy {
		return
	}

	if c.StatusIsActive(radianceNAICDKey) {
		return
	}

	c.AddStatus(radianceNAICDKey, 0.5*60, true)
	c.addRadiance(5)
}

func (c *char) skillInit() {
	c.Core.Events.Subscribe(event.OnLunarCrystallize, func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}
		if c.StatusIsActive(radianceLCrICDKey) {
			return
		}
		c.AddStatus(radianceLCrICDKey, 4*60, true)
		c.addRadiance(35)
	}, "zibai-radiance-lcr")
}

func (c *char) addRadiance(amt float64) {
	amt *= c.c6RadianceEff()
	c.radiance = min(c.radiance+amt, maxRadiance)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprint("Gained ", amt, " radiance (", c.radiance, ")"), glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) consumeRadiance(amt float64) {
	c.radiance = max(c.radiance-amt, 0)
	if c.Core.Flags.LogDebug {
		c.Core.Log.NewEvent(fmt.Sprint("Consumed ", amt, " radiance (", c.radiance, ")"), glog.LogCharacterEvent, c.Index())
	}
}
