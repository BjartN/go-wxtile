package example

import (
	"fmt"
	"github.com/bjartn/go-wxtile/tile"
	"github.com/bjartn/go-wxtile/data"
	"time"
	"image"
	"runtime"
	"github.com/llgcode/draw2d/draw2dimg"
	"image/color"
	"sync"
	"log"
	"strconv"
)

func QueryProjectAndResample() {

	// maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	world := tile.Bounds{tile.Point{-20000000,-20000000}, tile.Point{20000000,20000000}}
	sizeX := 1000
	sizeY:= 1000

	//get data
	grid := data.ParseJson("c:\\data\\grib-converted\\2015050600\\gfs.t00z.pgrb2.2p50.f000.2t.json")

	//check invariant
	if(grid.Header.ScanMode!=0){
		return
	}

	//draw latlon grid
	g := tile.ToRegularProjectedGrid(grid)
	i1 := tile.DrawGrid(&g)
	tile.SaveImage(i1,"c:\\temp\\lat-lon-grid.png")

	fmt.Println("")

	iParallel := tile.ResampleAndDrawConcurrent(g,world, sizeX,sizeY)
	tile.SaveImage(iParallel,"c:\\temp\\mercator-resampled-paralell-image.png")

	fmt.Println("")

	iSequential:=resampleAndDrawSequential(g,world, sizeX,sizeY)

	//prj.DrawX(iSequential)
	writeShapeToImage(iSequential,world)

	tile.SaveImage(iSequential,"c:\\temp\\mercator-resampled-sequential-grid.png")
	fmt.Println("Done")

	//var input string
	//fmt.Scanln(&input)

}


func ResampleConcurrent() {

	world := tile.Bounds{tile.Point{-20000000,-20000000}, tile.Point{20000000,20000000}}
	sizeX := 256
	sizeY:= 256
	grid := data.ParseJson("c:\\data\\grib-converted\\2015050600\\gfs.t00z.pgrb2.2p50.f000.2t.json")
	g := tile.ToRegularProjectedGrid(grid)


	N:=5000

	var wg sync.WaitGroup
	wg.Add(N)

	counter:=0
	start :=time.Now()
	f:= func(){
		defer wg.Done()
		i:=resampleAndDrawSequential(g,world, sizeX,sizeY)
		counter++
		tile.SaveImage(i, fmt.Sprintf("c:\\temp\\tiles\\" + strconv.Itoa(counter) + ".png"))
		//fmt.Print(".")
	}

	for i:=0; i<N; i++ {
		go f()
	}

	wg.Wait()

	elapsed := time.Since(start)

	perTile :=int( (elapsed.Seconds()*1000) / float64(N))
	log.Println("Per tile ", perTile, " ms")
	log.Printf( "%v concurrent tile operations took %s",N, elapsed)
	fmt.Println("Done")

}


func writeShapeToImage(m *image.RGBA,  b tile.Bounds){
	defer tile.TimeTrack(time.Now(), fmt.Sprint("WriteShapesToImage"))

	r,err := data.OpenShapeReader("c:\\dev\\intelli-go\\data\\pressure.shp")
	if(err!=nil){
		panic("I don't know how to handle errors")
	}

	width  := m.Bounds().Size().X
	height := m.Bounds().Size().X

	s1 := tile.Scale(b.Min.X,b.Max.X, 0, float64(width-1))
	s2 := tile.Scale(b.Max.Y,b.Min.Y, 0, float64(height-1))

	gc:=draw2dimg.NewGraphicContext(m)
	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{255, 255,255, 255})
	gc.SetLineWidth(1)

	for r.Next(){
		pixels := scaleMany(r.Polyline.Points, s1,s2)
		tile.DrawLine(gc, &tile.PolyLine{Points:pixels})
	}
}

func scaleMany ( p []tile.Point,scaleX func(x float64) float64 ,scaleY func(x float64) float64 ) []tile.Point {
	a:=make([]tile.Point, len(p))
	for i:= range p {
		a[i] = tile.Point{scaleX(p[i].X), scaleY(p[i].Y	)}
	}
	return a
}

func  resampleAndDrawSequential (g tile.Grid, b tile.Bounds, width int,height int) *image.RGBA {
	//defer prj.TimeTrack(time.Now(), fmt.Sprint("ResampleAndDrawSequential ", width, "*",height, "=", width*height/1000, "k"))

	g2 := tile.Resample(g, b, width, height)
	i2 := tile.DrawGrid(&g2)

	return i2
}


