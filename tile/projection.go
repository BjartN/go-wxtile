package tile

import (
	"math"
)

const MercatorPole = 20037508.34

type Projector interface {
	ToProjection(lon, lat float64) (Point)
	FromProjection(x, y float64) (Point)
}

type MercatorProjector struct {

}

//ToProjection will project a point into the Mercator projection.
func (MercatorProjector) ToProjection(lon, lat float64) (Point) {
	x := MercatorPole / 180.0 * lon
	y := math.Log(math.Tan((90.0 + lat) * math.Pi / 360.0)) / math.Pi * MercatorPole
	y = math.Max(-MercatorPole, math.Min(y, MercatorPole))

	if(x>=MercatorPole){
	 	x-=MercatorPole*2
	}

	return Point{x,y}
}

//FromProjection finds the geographic coordinates of a Mercator point.
func (MercatorProjector) FromProjection(x, y float64) (Point) {
	lon := x * 180.0 / MercatorPole
	lat := 180.0 / math.Pi * (2 * math.Atan(math.Exp((y / MercatorPole) * math.Pi)) - math.Pi / 2.0)

	if(lon<0){
		lon+=360
	}

	return Point{lon,lat}
}
