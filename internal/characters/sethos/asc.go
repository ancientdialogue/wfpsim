package sethos

// returns the amount of time to save, and the amount of energy considered
func (c *char) a1Calc() (int, float64) {
	if c.Base.Ascension < 1 {
		return 0, 0
	}
	energy := min(c.Energy, 20)
	// floor or round the skip?
	return int(0.285 * energy * 60), energy
}
