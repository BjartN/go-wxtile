package tile

import (
	"fmt"
	"math"
	"time"
)

//ToRegularProjectedGrid converts a regular geographic grid to a regular Mercator grid.
func ToRegularProjectedGrid(grid Grid) Grid {
	h := grid.Header
	projector := MercatorProjector{}

	defer TimeTrack(time.Now(), fmt.Sprint("Project ", h.Nx, "*", h.Ny, "=", (h.Nx*h.Ny)/1000, "k"))

	//coordinates on y axis
	projYs := make([]float64, h.Ny, h.Ny)
	for y := 0; y < h.Ny; y++ {
		lat := h.La1 - float64(y)*h.Dy
		lon := h.Lo1
		projYs[y] = projector.ToProjection(lon, lat).Y
	}

	//coordinates in x axis
	projXs := make([]float64, h.Nx, h.Nx)
	for x := 0; x < h.Nx; x++ {
		lat := h.La1
		lon := h.Lo1 + float64(x)*h.Dx
		projXs[x] = projector.ToProjection(lon, lat).X
	}

	//irregular projection grid
	pointValues := make([]PointValue, len(grid.Data))
	for y := 0; y < h.Ny; y++ {
		for x := 0; x < h.Nx; x++ {
			idx := y*h.Nx + x
			value := grid.Data[idx]

			pointValues[idx] = PointValue{Point{projXs[x], projYs[y]}, value}
		}
	}

	b := bounds(grid)
	projBounds := Bounds{Point{MercatorPole * -1, MercatorPole * -1}, Point{MercatorPole, MercatorPole}}
	dx := (MercatorPole * 2) / float64(h.Nx-1)
	dy := (MercatorPole * 2) / float64(h.Ny-1)

	result := make([]float32, len(grid.Data))

	//iterate grid points in regular projected grid, calculate lat/lon to find surrounding grid points in the original grid
	idx := 0
	for y := projBounds.Max.Y; y >= projBounds.Min.Y; y -= dy {
		for x := projBounds.Min.X; x <= projBounds.Max.X; x += dx {
			latlon := projector.FromProjection(x, y)

			idxX := int(((latlon.X - b.Min.X) / (b.Width())) * float64(h.Nx-1))
			idxY := int(((b.Max.Y - latlon.Y) / (b.height())) * float64(h.Ny-1))

			ul := pointValues[idxY*h.Nx+idxX]
			ll := pointValues[(idxY+1)*h.Nx+idxX]
			ur := pointValues[(idxY)*h.Nx+(idxX+1)]
			lr := pointValues[(idxY+1)*h.Nx+(idxX+1)]

			invariantOk := ul.X < ur.X && ll.Y < ul.Y && ll.X <= x && lr.X >= x && ll.Y <= y && ul.Y >= y

			if !invariantOk {
				fmt.Println("Invariant failed")
			}

			p := Point{x, y}
			v := BilinearInterpolation(&ll, &ul, &lr, &ur, &p)

			result[idx] = float32(v)

			idx++

		}
	}

	return Grid{Header: GridHeader{
		ScanMode: grid.Header.ScanMode,
		Nx:       grid.Header.Nx,
		Ny:       grid.Header.Ny,
		Dx:       dx,
		Dy:       dy,
		Lo1:      projBounds.Min.X,
		La1:      projBounds.Max.Y,
		Lo2:      projBounds.Max.X,
		La2:      projBounds.Min.Y}, Data: result}
}

func bounds(grid Grid) Bounds {

	p1 := Point{grid.Header.Lo1, grid.Header.La1}
	p2 := Point{grid.Header.Lo2, grid.Header.La2}

	xMin := math.Min(p1.X, p2.X)
	xMax := math.Max(p1.X, p2.X)
	yMin := math.Min(p1.Y, p2.Y)
	yMax := math.Max(p1.Y, p2.Y)

	return Bounds{Point{xMin, yMin}, Point{xMax, yMax}}
}

//Create function for scaling float64 -> float64
func Scale(from1 float64, from2 float64, to1 float64, to2 float64) func(x float64) float64 {

	length1 := from2 - from1
	length2 := to2 - to1

	return func(x float64) float64 {
		ratio := (x - from1) / length1

		return to1 + ratio*length2
	}
}
