package xiao

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

// Normal attack damage queue generator
// relatively standard with no major differences versus other characters
func (c *char) Attack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionAttack, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		AttackTag:  core.AttackTagNormal,
		ICDTag:     core.ICDTagNormalAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
	}

	for i, mult := range attack[c.NormalCounter] {
		ai.Mult = mult[c.TalentLvlAttack()]
		c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-5+i, f-5+i)
	}

	c.AdvanceNormalIndex()

	// return animation cd
	return f, a
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
// Note that his CAs share an ICD with his NAs when he is under the effects of his burst
// TODO: No information available on whether regular CAs follow a similar pattern
func (c *char) ChargeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionCharge, p)

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge",
		AttackTag:  core.AttackTagExtra,
		ICDTag:     core.ICDTagExtraAttack,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), f-1, f-1)

	//return animation cd
	return f, a
}

// Plunge normal falling attack damage queue generator
// Standard - Always part of high/low plunge attacks
func (c *char) PlungeAttack(delay int) (int, int) {
	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Plunge (Normal)",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Physical,
		Durability: 25,
		Mult:       plunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(0.1, false, core.TargettableEnemy), delay, delay)

	//return animation cd
	return delay, delay
}

// High Plunge attack damage queue generator
// Use the "plunge_hits" optional argument to determine how many plunge falling hits you do on the way down
// Default = 0
func (c *char) HighPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionHighPlunge, p)

	plungeHits, ok := p["plunge_hits"]
	if !ok {
		plungeHits = 0 // Number of normal plunge hits
	}

	for i := 0; i < plungeHits; i++ {
		// Add plunge attack in each frame leading up to final hit for now - not sure we have clear mechanics on this
		// TODO: Perhaps amend later, but functionally in combat you usually get at most one of these anyway
		c.PlungeAttack(f - i - 2)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "High Plunge",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       highplunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-1, f-1)

	//return animation cd
	return f, a
}

// Low Plunge attack damage queue generator
// Use the "plunge_hits" optional argument to determine how many plunge falling hits you do on the way down
// Default = 0
func (c *char) LowPlungeAttack(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionLowPlunge, p)

	plungeHits, ok := p["plunge_hits"]
	if !ok {
		plungeHits = 0 // Number of normal plunge hits
	}

	for i := 0; i < plungeHits; i++ {
		// Add plunge attack in each frame leading up to final hit for now - not sure we have clear mechanics on this
		// TODO: Perhaps amend later, but functionally in combat you usually get at most one of these anyway
		c.PlungeAttack(f - i - 2)
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Low Plunge",
		AttackTag:  core.AttackTagPlunge,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		StrikeType: core.StrikeTypeBlunt,
		Element:    core.Physical,
		Durability: 25,
		Mult:       lowplunge[c.TalentLvlAttack()],
	}
	c.Core.Combat.QueueAttack(ai, core.NewDefCircHit(2, false, core.TargettableEnemy), f-1, f-1)

	//return animation cd
	return f, a
}

// Skill attack damage queue generator
// Additionally implements A4
// Using Lemniscatic Wind Cycling increases the DMG of subsequent uses of Lemniscatic Wind Cycling by 15%. This effect lasts for 7s and has a maximum of 3 stacks. Gaining a new stack refreshes the duration of this effect.
func (c *char) Skill(p map[string]int) (int, int) {

	f, a := c.ActionFrames(core.ActionSkill, p)
	// According to KQM library, double E is only 60 frames long whereas single cast is 36
	// No idea how this works, but add a special case to reduce the frames of the 2nd cast
	// TODO: No data listed on how 3 casts work - this might be too few frames compared to actual 3x usage
	if c.Core.LastAction.Target == core.Xiao && c.Core.LastAction.Typ == core.ActionSkill {
		f = 60 - f
		a = 60 - a
	}

	// Add damage based on A4
	if c.a4Expiry <= c.Core.F {
		c.Tags["a4"] = 0
	}

	ai := core.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lemniscatic Wind Cycling",
		AttackTag:  core.AttackTagElementalArt,
		ICDTag:     core.ICDTagElementalArt,
		ICDGroup:   core.ICDGroupDefault,
		Element:    core.Anemo,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}
	snap := c.Snapshot(&ai)
	c.Core.Combat.QueueAttackWithSnap(ai, snap, core.NewDefCircHit(2, false, core.TargettableEnemy), f)

	// Text is not explicit, but assume that gaining a stack while at max still refreshes duration
	c.Tags["a4"]++
	c.a4Expiry = c.Core.F + 420
	if c.Tags["a4"] > 3 {
		c.Tags["a4"] = 3
	}

	// Cannot create energy during burst uptime
	if c.Core.Status.Duration("xiaoburst") > 0 {
	} else {
		c.QueueParticle("xiao", 3, core.Anemo, f+100)
	}

	// C6 handling - can use skill ignoring CD and without draining charges
	// Can simply return early
	if c.Base.Cons == 6 && c.Core.Status.Duration("xiaoc6") > 0 {
		c.Core.Log.NewEvent("xiao c6 active, Xiao E used, no charge used, no CD", core.LogCharacterEvent, c.Index, "c6 remaining duration", c.Core.Status.Duration("xiaoc6"))
		return f, a
	}

	c.SetCD(core.ActionSkill, 600)

	return f, a
}

// Sets Xiao's burst damage state
func (c *char) Burst(p map[string]int) (int, int) {
	f, a := c.ActionFrames(core.ActionBurst, p)
	var HPicd int
	var HPlost float64
	HPicd = 0

	// Per previous code, believe that the burst duration starts ticking down from after the animation is done
	// TODO: No indication of that in library though
	c.Core.Status.AddStatus("xiaoburst", 900+f)
	c.qStarted = c.Core.F

	// HP Drain - removes HP every 1 second tick after burst is activated
	// Per gameplay video, HP ticks start after animation is finished
	for i := f + 60; i < 900+f; i++ {
		c.AddTask(func() {
			if c.Core.Status.Duration("xiaoburst") > 0 && c.Core.F >= HPicd {
				HPlost = burstDrain[c.TalentLvlBurst()] * c.HPCurrent
				c.HPCurrent = c.HPCurrent - HPlost
				c.Core.Log.NewEvent("xiao hp drain", core.LogCharacterEvent, c.Index, "current HP", c.HPCurrent, "HP lost", HPlost)
				c.Core.Events.Emit(core.OnCharacterHurt, HPlost)
				HPicd = c.Core.F + 60
			}
		}, "xiaoburst-hp-drain", i)
	}

	// Checked gameplay - burst starts ticking down from activation. CD is 16.6 seconds after animation is done
	c.SetCDWithDelay(core.ActionBurst, 18*60, 39)
	c.ConsumeEnergy(39)

	return f, a
}
