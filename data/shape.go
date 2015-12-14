package data

import (
	"github.com/bjartn/go-wxtile/tile"
	"github.com/jonas-p/go-shp"
	"fmt"
	"time"
)

type ShapeReader struct {
	projector tile.Projector
	reader    shp.Reader
	err       error
	Polyline  tile.PolyLine
}

func (sr *ShapeReader) Close() {
	sr.reader.Close()
}

func (sr *ShapeReader) Next() bool {

	for sr.reader.Next() {
		_, p := sr.reader.Shape()

		switch p.(type) {
		case *shp.PolyLine:
			polyLine := p.(*shp.PolyLine)
			ph := tile.PolyLine{};
			ph.Points = make([]tile.Point, len(polyLine.Points))

			for i := range polyLine.Points {
				p:=polyLine.Points[i]
				curr := sr.projector.ToProjection(p.X, p.Y)

				if(p.X>180) { //fix only for that one strange shape file I had that time
					curr.X = tile.MercatorPole
				}

				ph.Points[i] = curr;
			}
			sr.Polyline = ph //concurrency issues here?
			return true
		}
	}

	return false
}

func OpenShapeReader(file string) (*ShapeReader, error) {
	defer tile.TimeTrack(time.Now(), fmt.Sprint("GetShape"))

	// open a shape file for reading
	shape, err := shp.Open(file)
	if err != nil {
		return nil, err
	}

	sr := ShapeReader{reader:*shape, projector:tile.MercatorProjector{}}

	return &sr,nil

}
