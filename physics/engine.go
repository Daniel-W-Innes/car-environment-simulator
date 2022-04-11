package physics

import (
	"math"
)

func toRadians(deg int) float64 {
	return float64(deg) * (math.Pi / 180.0)
}
func xComponent(d float64, angle int) float64 {
	if angle < 90 {
		return d * math.Sin(toRadians(angle))
	}
	if angle < 180 {
		return d * math.Cos(toRadians(angle-90))
	}
	if angle < 270 {
		return -d * math.Sin(toRadians(angle-180))
	}
	return -d * math.Cos(toRadians(angle-270))
}

func yComponent(d float64, angle int) float64 {
	if angle < 90 {
		return d * math.Cos(toRadians(angle))
	}
	if angle < 180 {
		return -d * math.Sin(toRadians(angle-90))
	}
	if angle < 270 {
		return -d * math.Cos(toRadians(angle-180))
	}
	return d * math.Sin(toRadians(angle-270))
}

func deltaD(dt, v, a, j float64) float64 {
	output := 0.0
	if v != 0 {
		output += v * dt
	}
	if a != 0 {
		output += 0.5 * a * math.Pow(dt, 2)
	}
	if j != 0 {
		output += 1.0 / 6 * j * math.Pow(math.Min(-a/j, dt), 3)
	}
	return output
}

func deltaV(dt, a, j float64) float64 {
	output := 0.0
	if a != 0 {
		output += a * dt
	}
	if j != 0 {
		output += 0.5 * j * math.Pow(math.Min(-a/j, dt), 2)
	}
	return output
}
func deltaA(dt, a, j float64) float64 {
	if j != 0 {
		return j * math.Min(-a/j, dt)
	}
	return 0
}

func deltaJ(dt, j float64) float64 {
	if j == 0 {
		return j
	}
	if j > 0 {
		return j - maxJ*dt
	} else {
		return j + maxJ*dt
	}
}
