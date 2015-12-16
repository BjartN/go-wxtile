package main
import (
	"github.com/bjartn/go-wxtile/tile"
	"github.com/bjartn/go-wxtile/data"
)


func main() {
	targetFolder:="c:\\temp\\tiles"

	maxZoom := 5

	//get data (grib file serialized using grib2json)
	grid := data.ParseJson("C:\\dev\\go-maps\\data\\gfs.t00z.pgrb2.2p50.f000.2t.json")

	//project to mercator
	g := tile.ToRegularProjectedGrid(grid)

	//make tiles
	tile.Tiler(g, maxZoom, targetFolder)

}




