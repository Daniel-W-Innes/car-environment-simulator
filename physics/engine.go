package physics

import (
	"math"
)

func toRadians(deg uint16) float64 {
	return float64(deg) * (math.Pi / 180.0)
}
func xComponent(d float64, angle uint16) float64 {
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

func yComponent(d float64, angle uint16) float64 {
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
	return v*dt + 0.5*a*math.Pow(dt, 2) + 1.0/6*j*math.Pow(math.Min(-a/j, dt), 3)
}

func deltaV(dt, a, j float64) float64 {
	return a*dt + 0.5*j*math.Pow(math.Min(-a/j, dt), 2)
}
func deltaA(dt, a, j float64) float64 {
	return j * math.Min(-a/j, dt)
}
