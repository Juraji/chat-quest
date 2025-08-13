package util

// MaxFloat32 returns the larger of two float32 values
func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
