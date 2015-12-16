package tile

import (
	"math"
	"time"
	"fmt"
	"strconv"
	"sync"
)

func Tiler(g Grid, maxZoom int, targetFolder string){
	defer TimeTrack(time.Now(), fmt.Sprint("Tiler with zoom ", maxZoom))

	tileSize := 256
	world := Bounds{
		Point{MercatorPole*-1,MercatorPole*-1},
		Point{MercatorPole,MercatorPole}}

	var wg sync.WaitGroup

	for z:=1; z<=maxZoom; z++ {
		sizef := math.Pow(2,float64(z))
		for x:=0; x<int(sizef); x++ {
			xf := float64(x)
			for y:=0; y<int(sizef); y++ {
				wg.Add(1)

				width := world.Width() / sizef
				yf := float64(y)

				//make bounds
				b := Bounds{
					Point{MercatorPole*-1 +  width * xf, MercatorPole - width * (yf + 1)},
					Point{MercatorPole*-1 +  width * (xf + 1), MercatorPole - width * yf}}

				fileName := fmt.Sprintf(targetFolder + "\\" + strconv.Itoa(int(z)) + "_" + strconv.Itoa(int(x)) + "_" + strconv.Itoa(int(y)) + ".png")

				go (func(b Bounds, fileName string) {

					//make tile
					gTile := Resample(g, b, tileSize, tileSize)

					//make image
					img := DrawGrid(&gTile)

					//save image to disc
					SaveImage(img, fileName)

					wg.Done()
				})(b,fileName)
			}
		}
	}

	wg.Wait()

}