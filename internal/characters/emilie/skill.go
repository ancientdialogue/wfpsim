package emilie

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
	tickInterval       = 120 // assume consistent 120f tick rate
	skillDuration      = 22 * 60
	particleICDKey     = "emilie-particle-icd"
	skillKey           = "emilie-skill"
	scentICD           = 115
	scentICDKey        = "scent-generation-icd"
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
	// always trigger electro no ICD on initial summon
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lumidouce Case (Summon)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}

	radius := 4.5
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 1.5}, radius),
		skillLumiSpawn,
		skillLumiHitmark,
	)

	// CD Delay is 18 frames, but things break if Delay > CanQueueAfter
	// so we add 18 to the duration instead. this probably mess up CDR stuff
	c.SetCD(action.ActionSkill, 14*60+skillLumiSpawn)

	c.queueLumiSpawn(skillLumiSpawn, skillLumiFirstTick, 1)

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

func (c *char) queueLumiSpawn(lumiSpawn, firstTick, lumidouceLvl int) {
	// calculate duration
	if lumidouceLvl < 1 || lumidouceLvl > 2 {
		return
	}
	if c.lumidouceLvl > 2 {
		return
	}
	c.Core.Tasks.Add(func() {
		c.lumidouceScents = 0
		c.lumidouceSrc = c.Core.F
		c.lumidouceScentLast = c.Core.F
		c.lumidouceLvl = lumidouceLvl
		c.Core.Tasks.Add(c.removeLumi(c.Core.F), skillDuration)
		player := c.Core.Combat.Player()
		c.lumidoucePos = geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 1.5}, player.Direction())
		c.AddStatus(skillKey, skillDuration, false)
		c.QueueCharTask(c.lumiTick(c.Core.F), firstTick)
		c.QueueCharTask(c.scentCheck(c.Core.F), 0)
		c.QueueCharTask(c.scentCollect(c.Core.F), 0)
		c.Core.Log.NewEvent("Lumidouce case activated", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+tickInterval)
	}, lumiSpawn)
}

func (c *char) lumiTick(src int) func() {
	return func() {
		// if src != lumidouceSrc then this is no longer the same lumidouce case, do nothing
		if src != c.lumidouceSrc {
			return
		}
		if !c.StatusIsActive(skillKey) {
			return
		}
		c.Core.Log.NewEvent("Lumidouce Case ticked", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+tickInterval).
			Write("src", src)
		// trigger damage
		switch c.lumidouceLvl {
		case 1:
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lumidouce Case (Level 1)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupEmilieLumidouce,
				StrikeType: attacks.StrikeTypePierce,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       skillLumidouce[0][c.TalentLvlSkill()],
			}
			ap := combat.NewCircleHit(
				c.lumidoucePos,
				c.Core.Combat.PrimaryTarget(),
				nil,
				1,
			)
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
		case 2:
			ai := combat.AttackInfo{
				ActorIndex: c.Index,
				Abil:       "Lumidouce Case (Level 2)",
				AttackTag:  attacks.AttackTagElementalArt,
				ICDTag:     attacks.ICDTagElementalArt,
				ICDGroup:   attacks.ICDGroupEmilieLumidouce,
				StrikeType: attacks.StrikeTypePierce,
				Element:    attributes.Dendro,
				Durability: 25,
				Mult:       skillLumidouce[1][c.TalentLvlSkill()],
			}
			ap := combat.NewCircleHit(
				c.lumidoucePos,
				c.Core.Combat.PrimaryTarget(),
				nil,
				1,
			)
			c.Core.QueueAttack(ai, ap, 0, 0, c.particleCB)
			c.Core.QueueAttack(ai, ap, 5, 5, c.particleCB)
		}

		// queue up next hit
		c.Core.Tasks.Add(c.lumiTick(src), tickInterval)
	}
}

func (c *char) removeLumi(src int) func() {
	return func() {
		// if src != lumidouceSrc then this is no longer the same lumidouce, do nothing
		if c.lumidouceSrc != src {
			c.Core.Log.NewEvent("Lumidouce Case not removed, src changed", glog.LogCharacterEvent, c.Index).
				Write("src", src)
			return
		}
		c.lumidouceLvl = -1
		c.lumidouceScentLast = -1
		c.lumidouceScents = -1
		c.lumidouceSrc = -1
		c.DeleteStatus(skillKey)
		c.Core.Log.NewEvent("Lumidouce Case removed", glog.LogCharacterEvent, c.Index).
			Write("src", src)
	}
}

func (c *char) scentCheck(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}

		c.QueueCharTask(c.scentCheck(src), 30)

		if c.StatusIsActive(scentICDKey) {
			return
		}

		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.lumidoucePos, nil, 8), nil)
		for _, v := range enemies {
			e, ok := v.(*enemy.Enemy)
			if !ok {
				continue
			}
			if e.IsBurning() {
				c.scentGen()
				c.QueueCharTask(c.scentCheck(src), 120)
				break
			}
		}
	}
}

func (c *char) scentGen() {
	// it's a gadget doing the scent collection so not hitlag affected?
	c.AddStatus(scentICDKey, scentICD, false)

	// TODO: Do scents expire??
	// Does this need to become a gadget?
	// TODO: What is the scent collection delay?
	c.Core.Tasks.Add(func() {
		c.scentCount++
		c.Core.Log.NewEvent("Scent generated", glog.LogCharacterEvent, c.Index).
			Write("total scents", c.scentCount)
	}, 20)
}

func (c *char) scentCollect(src int) func() {
	return func() {
		if c.lumidouceSrc != src {
			return
		}

		c.QueueCharTask(c.scentCollect(src), 30)

		if c.scentCount > 0 {
			c.scentCount--
			c.lumidouceScents++
			c.lumidouceScentLast = c.Core.F
			c.Core.Log.NewEvent("Scent Collected", glog.LogCharacterEvent, c.Index).
				Write("total scents left", c.scentCount).
				Write("total scents collected", c.lumidouceScents)
			if c.lumidouceLvl == 2 {
				c.a1()
			}

			if c.lumidouceScents > 2 {
				c.lumidouceLvl = 2
				c.Core.Log.NewEvent("Lumidouce upgraded", glog.LogCharacterEvent, c.Index)
			}
			return
		}
		if c.Core.F-c.lumidouceScentLast >= 8*60 {
			c.lumidouceLvl = 1
		}
	}
}
