package chasca

var c2icd = "chasca-c2-icd"

func (c *char) c2stacks() int {
	if c.Base.Cons > 2 {
		return 1
	}
	return 0
}
