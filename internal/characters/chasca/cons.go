package chasca

var c2icd = "chasca-c2-icd"
var c4energy = "chasca-c4-energy"
var c4icd = "chasca-c4-icd"

func (c *char) c2stacks() int {
	if c.Base.Cons > 2 {
		return 1
	}
	return 0
}

func (c *char) c4energy() {
	if c.Base.Cons < 4 {
		return
	}
	if !c.StatusIsActive(c4energy) {
		c.AddEnergy(c4energy, 1.5)
	}

	c.AddStatus(c4energy, -1, false)
}
