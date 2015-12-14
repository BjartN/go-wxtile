package tile

import (
	"testing"
	"github.com/bjartn/go-wxtile/tile"
	"fmt"
)

func TestShouldProject(t *testing.T){
	m := tile.MercatorProjector{}

	/*
		180.0625 -50.869140625
		179.9375 -50.831640625000006
		-2.0030550871826388e+07 -6.598178774984599e+06
		2.003055087182639e+07 -6.59156675680389e+06
	*/

	lng:=180.0625
	lat:=-50.869140625
	p:= m.ToProjection(lng,lat)
	fmt.Println(p.X, p.Y)

	lng2:=179.9375
	lat2:=-50.831640625000006
	p2:= m.ToProjection(lng2,lat2)
	fmt.Println(p2.X, p2.Y)
}