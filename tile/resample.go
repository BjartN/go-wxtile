package tile

import (
	"math"
)

//Resample resamples the specified area of a regular grid.
func Resample(grid Grid, b Bounds, width int, height int) Grid {
	//defer TimeTrack(time.Now(), fmt.Sprint("Resample ", width, "*",height, "=", width*height/1000, "k"))
	h := grid.Header;

	//xMin:xMax scales to 0:width-1
	xMin := (b.Min.X-h.Lo1) / h.Dx
	xMax  := (b.Max.X-h.Lo1) / h.Dx

	//yMax:yMin scales to 0:height-1
	yMin  := (h.La1 - b.Max.Y) / h.Dy;
	yMax := (h.La1- b.Min.Y) / h.Dy;

	dx := (xMax-xMin) / float64(width)
	dy := (yMax-yMin) / float64(height)

	yMap :=   make([]int,height)
	xMap :=   make([]int,width)
	yWeight :=   make([]float64,height)
	xWeight :=   make([]float64,width)

	for y := 0; y < height-1; y++ {
		idx :=  yMin + float64(y)*dy
		yMap[y] = int(idx)
		yWeight[y] = idx - math.Floor(idx)
	}

	for x := 0; x < width-1; x++ {
		idx:=xMin + float64(x)*dx
		xMap[x] = int(idx)
		xWeight[x] = idx - math.Floor(idx)
	}

	newHeader := GridHeader{
		Nx : width,
		Ny : height,
		Dx : 1,
		Dy : 1,
		ScanMode:0,
		Lo1 : 0,
		La1 : float64(height)-1,
		Lo2 : float64(width)-1,
		La2 : 0}

	newGrid :=Grid{make([]float32,width*height),newHeader}

	for y := 0; y < height-1; y++ {
		wy :=yWeight[y]
		origY := yMap[y]

		for x := 0; x < width-1; x++ {
			wx :=xWeight[x]
			origX := xMap[x]

			ul := grid.GetValueAtIdx(origX, origY)
			ur := grid.GetValueAtIdx(origX+1, origY)
			ll := grid.GetValueAtIdx(origX, origY+1)
			lr := grid.GetValueAtIdx(origX+1, origY+1)

			xUpperValue := float64(ul.Value)  * (1-wx) + float64(ur.Value) * wx
			xLowerValue := float64(ll.Value)  * (1-wx) + float64(lr.Value) * wx

			v := xUpperValue * (1-wy) + xLowerValue * wy

			newGrid.Data[newGrid.GetIdx(x,y)] =float32(v)
		}
	}


	return newGrid
}