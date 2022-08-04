package internal

type CelestialBody struct {
}

func (CelestialBody) Type() string {
	return "CelestialBody"
}

type Sphere struct {
	Radius float64
}

func (Sphere) Type() string {
	return "Sphere"
}

type Position struct {
	X float64
	Y float64
}

func (Position) Type() string {
	return "Position"
}
