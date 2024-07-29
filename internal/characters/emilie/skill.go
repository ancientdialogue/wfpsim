package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillLumiSpawn     = 18 // same as CD start
	skillLumiHitmark   = 38
	skillLumiFirstTick = 64
	tickInterval       = 90 // assume consistent 59f tick rate
	particleICDKey     = "emilie-particle-icd"
	scentICDKey        = "emilie-scent-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 16
	skillFrames[action.ActionSwap] = 42

	skillRecastFrames = frames.InitAbilSlice(37)
	skillRecastFrames[action.ActionAttack] = 36
	skillRecastFrames[action.ActionBurst] = 35
	skillRecastFrames[action.ActionDash] = 4
	skillRecastFrames[action.ActionJump] = 5
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lumidouce Case (Summon)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupEmilieLumidouce,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 3}, 4.5),
		skillLumiSpawn,
		skillLumiHitmark,
	)

	player := c.Core.Combat.Player()
	c.lumidouceLvl = 0
	c.lumidouceSrc = c.Core.F
	c.lumidoucePos = geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 1.5}, player.Direction())

	if !c.lumidouceCheck {
		c.lumidouceCheck = true
		c.checkScents()
	}

	c.Core.Tasks.Add(c.lumiTick(c.Core.F), skillLumiFirstTick)
	c.Core.Tasks.Add(c.removeLumi(c.Core.F), 22*60)

	c.SetCD(action.ActionSkill, 14*60+skillLumiSpawn)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 2.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Dendro, c.ParticleDelay)
}

func (c *char) lumiTick(src int) func() {
	return func() {
		if src != c.lumidouceSrc {
			return
		}

		ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 1)

		if c.lumidouceLvl >= 2 {
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lumidouce Case Lv2",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupEmilieLumidouce,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       skillLumidouce[1][c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
			c.Core.QueueAttack(ai, ap, 0, 10, c.particleCB)
		} else {
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lumidouce Case Lv1",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupEmilieLumidouce,
				StrikeType: attacks.StrikeTypeDefault,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       skillLumidouce[0][c.TalentLvlSkill()],
			}
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
		}

		c.Core.Tasks.Add(c.lumiTick(src), tickInterval)
	}
}

func (c *char) checkScents() {
	if c.lumidouceSrc == -1 {
		return
	}
	isBurning := false
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 20), nil)
	for _, v := range enemies {
		e, ok := v.(*enemy.Enemy)
		if !ok {
			continue
		}
		if e.IsBurning() {
			isBurning = true
			break
		}
	}
	if isBurning && !c.StatModIsActive(scentICDKey) {
		c.AddStatus(scentICDKey, 2*60, false)
		c.genScents()
	}
	c.QueueCharTask(c.checkScents, 30)
}

func (c *char) genScents() {
	if c.lumidouceLvl < 4 {
		c.lumidouceLvl++
	}
	if c.lumidouceLvl == 4 {
		c.lumidouceLvl = 2
		c.a1()
	}
}

func (c *char) removeLumi(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}
		c.lumidouceCheck = false
		c.lumidouceSrc = -1
		c.lumidouceLvl = 0
	}
}
