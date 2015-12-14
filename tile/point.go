package tile


type Point struct {
	X float64
	Y float64

}

type PointValue struct {
	Point
	Value float32
}

type PolyLine struct {
	Bounds
	Points    []Point
}

type Bounds struct {
	Min Point
	Max Point
}

func (b Bounds) Width() float64 {
	return b.Max.X - b.Min.X
}

func (b Bounds) height() float64 {
	return b.Max.Y - b.Min.Y
}

