package tile

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"image/png"
	"time"
	"sync"
	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	blue color.Color = color.RGBA{0, 0, 255, 255}
	palette []color.Color= []color.Color{
		color.RGBA{183, 219, 243, 255},
		color.RGBA{162, 201, 227, 255},
		color.RGBA{140, 183, 212, 255},
		color.RGBA{119, 165, 196, 255},
		color.RGBA{97, 148, 180, 255},
		color.RGBA{76, 130, 164, 255},
		color.RGBA{54, 112, 149, 255},
		color.RGBA{33, 94, 133, 255},
		color.RGBA{133, 33, 33, 255},
		color.RGBA{154, 27, 29, 255},
		color.RGBA{175, 22, 25, 255},
		color.RGBA{195, 16, 22, 255},
		color.RGBA{216, 11, 18, 255},
		color.RGBA{237, 5, 14, 255},
		color.RGBA{242, 52, 11, 255},
		color.RGBA{246, 99, 7, 255},
		color.RGBA{251, 145, 4, 255},
		color.RGBA{255, 192, 0, 255},
})

func DrawX(img *image.RGBA) {
	defer TimeTrack(time.Now(), fmt.Sprint("DrawX "))

	gc:=draw2dimg.NewGraphicContext(img)

	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(5)

	gc.MoveTo(0,0)
	gc.LineTo(float64(img.Bounds().Size().X-1), float64(img.Bounds().Size().Y-1))

	gc.MoveTo(0,float64(img.Bounds().Size().Y-1))
	gc.LineTo(float64(img.Bounds().Size().X-1), 0)
	gc.Stroke()

	//gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()

}

func SaveImage( m *image.RGBA, file string) {
	defer TimeTrack(time.Now(), fmt.Sprint("Encode ", m.Rect.Size().X, "*", m.Rect.Size().Y, "=", (m.Rect.Size().X * m.Rect.Size().Y)/1000, "k"))

	w, _ := os.Create(file)
	defer w.Close()
	png.Encode(w, m)
}

func DrawLine(gc *draw2dimg.GraphicContext, line *PolyLine ){

	if len(line.Points)<2{
		return
	}

	gc.MoveTo(line.Points[0].X,line.Points[0].Y)

	for i:= range line.Points {
		gc.LineTo(line.Points[i].X,line.Points[i].Y)
	}
	gc.Stroke()
	//gc.Close()
}

func DrawGridInto(grid *Grid,m *image.RGBA, offsetX int, offsetY int) {

	max := 50.0
	min := -20.0

	scale := Scale(min, max, 0, float64(len(palette)-1))

	colorMapOffset := 20.0
	colorMap:= make([]color.Color, int(max-min)+1)

	for v := min; v<=max; v++ {
		colorMap[int(v+colorMapOffset)] = palette[int(scale(v))]
	}

	for y := 0; y < grid.Header.Ny; y++ {
		for x := 0; x < grid.Header.Nx; x++ {
			value := float64(grid.Data[grid.GetIdx(x,y)] - 273.15)

			if value==-9999{
				continue
			}
			if value<min{
				value = min
			}
			if value>max {
				value=max
			}

			m.Set(x+offsetX, y+offsetY, colorMap[int(colorMapOffset+value)])
		}
	}

}

func DrawGrid(grid *Grid) *image.RGBA{
	m := image.NewRGBA(image.Rect(0, 0, grid.Header.Nx, grid.Header.Ny))
	DrawGridInto(grid,m,0,0)
	return m
}

func DrawGridTest(grid Grid) *image.RGBA{
	h:=grid.Header
	defer TimeTrack(time.Now(), fmt.Sprint("Paint Simple ", h.Nx, "*", h.Ny, "=", (h.Nx * h.Ny)/1000, "k"))

	m := image.NewRGBA(image.Rect(0, 0, grid.Header.Nx, grid.Header.Ny))

	for y := 0; y < grid.Header.Ny; y++ {
		for x := 0; x < grid.Header.Nx; x++ {
			m.Set(x, y, blue)
		}
	}

	return m;
}


func ResampleAndDrawConcurrent (g Grid, b Bounds, width int,height int) *image.RGBA {
	defer TimeTrack(time.Now(), fmt.Sprint("ResampleAndDrawConcurrent ", width, "*",height, "=", width*height/1000, "k"))

	/*
		1) Break up grid into separate chunks by vertical slices
		2) Resample and draw into image in parallel
	 */

	//Break up into vertical slizes (but not to small)
	N := 4
	w := width  /N
	wCoord := b.Width() / float64(N)
	grids := make([]Grid, N);

	a := make([]Bounds, N)
	for i:=0; i<N; i++ {
		xMin := b.Min.X + wCoord*float64(i)
		a[i] = Bounds{Point{xMin,b.Min.Y}, Point{xMin+wCoord,b.Max.Y}}
	}

	//Resample and paint in paralell
	var wg sync.WaitGroup
	wg.Add(len(grids))

	 image.NewPaletted(image.Rect(0, 0, width,height), palette )

	img := image.NewRGBA(image.Rect(0, 0, width,height))


	for i,_ := range a {
		go func (i int) {
			defer wg.Done()
			gSample := Resample(g, a[i], w,height)
			DrawGridInto(&gSample,img, i*w ,0)

		} (i);
	}

	wg.Wait()
	return img
}