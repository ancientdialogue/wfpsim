package emilie

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	skillFrames       []int
	skillRecastFrames []int
)

const (
	skillLumiSpawn     = 18 // same as CD start
	skillLumiHitmark   = 38
	skillLumiFirstTick = 64
	tickInterval       = 120 // assume consistent 59f tick rate
	particleICDKey     = "emilie-particle-icd"
	skillKey           = "emilie-skill"
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
		ICDGroup:   attacks.ICDGroupEmilieLumidouce,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillDMG[c.TalentLvlSkill()],
	}

	radius := 2.0
	// hitmark is 5 frames after oz spawns
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{Y: 1.5}, radius),
		skillLumiSpawn,
		skillLumiHitmark,
	)

	// CD Delay is 18 frames, but things break if Delay > CanQueueAfter
	// so we add 18 to the duration instead. this probably mess up CDR stuff
	c.SetCD(action.ActionSkill, 14*60+skillLumiSpawn)

	// set on field oz to be this one
	c.queueLumi(skillLumiSpawn, skillLumiFirstTick)

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
	c.AddStatus(particleICDKey, 0.1*60, true)
	if c.Core.Rand.Float64() < .67 {
		// TODO: this delay used to be 120
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
	}
}

func (c *char) queueLumi(lumiSpawn, firstTick int) {
	// calculate oz duration
	dur := 22 * 60
	spawnFn := func() {
		// setup variables for tracking oz
		c.lumidouceSrc = c.Core.F
		// queue up oz removal at the end of the duration for gcsl conditional
		c.Core.Tasks.Add(c.removeLumi(c.Core.F), dur)

		player := c.Core.Combat.Player()
		c.lumidoucePos = geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 1.5}, player.Direction())

		c.Core.Tasks.Add(c.lumiTick(c.Core.F), firstTick)
		c.Core.Log.NewEvent("Oz activated", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+tickInterval)
	}
	if lumiSpawn > 0 {
		c.Core.Tasks.Add(spawnFn, lumiSpawn)
		return
	}
	spawnFn()
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
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       fmt.Sprintf("Lumidouce Case (%v)", src),
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagElementalArt,
			ICDGroup:   attacks.ICDGroupFischl,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skillLumidouce[0][c.TalentLvlSkill()],
		}
		ap := combat.NewBoxHit(
			c.lumidoucePos,
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		)
		c.Core.QueueAttack(ai, ap, 0, 0)

		// queue up next hit only if next hit oz is still active
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
		c.Core.Log.NewEvent("Lumidouce Case removed", glog.LogCharacterEvent, c.Index).
			Write("src", src)
	}
}
