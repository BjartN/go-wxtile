package tile

type Grid struct {
	Data   []float32
	Header GridHeader
}

type GridHeader struct {
	ScanMode int
	Nx       int
	Ny       int
	Dx       float64
	Dy       float64
	Lo1      float64
	Lo2      float64
	La1      float64
	La2      float64
}


func (g *Grid) Width() float64 {
	return g.Header.Lo2 - g.Header.Lo1
}

func (g *Grid) Height() float64 {
	return g.Header.La1 - g.Header.La2
}

func (g *Grid) GetIdx(idxX int, idxY int) int {
	return idxY * g.Header.Nx + idxX
}

//Get value at a given index, assuming an y-axis pointing north-to-south
func (g *Grid) GetValueAtIdx(idxX int, idxY int) PointValue{
	h:= g.Header

	if(idxX<0 || idxX>=h.Nx) || (idxY<0 || idxY>=h.Ny){
		//fmt.Println("Out of bounds")
		return PointValue{Point{h.Lo1 + float64(idxX)*h.Dx, h.La1 - float64(idxY)*h.Dy }, float32(-9999)}
	}

	v:= g.Data[idxY * h.Nx + idxX]

	x:= h.Lo1 + float64(idxX)*h.Dx
	y:= h.La1 - float64(idxY)*h.Dy

	return PointValue{Point{x, y },v}
}


//Use bilinear interpolation to get the value at a point within the grid
func (g Grid) GetValueAt(p *Point) float32 {

	if(p.X<g.Header.Lo1 || p.X>g.Header.Lo2){
		return float32(-9999);
	}

	if(p.Y>g.Header.La1 || p.Y<g.Header.La2){
		return float32(-9999);
	}

	idxX :=  int(((p.X - g.Header.Lo1) / g.Width()) * float64(g.Header.Nx-1))
	idxY :=  int(((g.Header.La1 - p.Y) / g.Height()) * float64(g.Header.Ny-1))

	ul := g.GetValueAtIdx(idxX, idxY)
	ur := g.GetValueAtIdx(idxX+1, idxY)
	ll := g.GetValueAtIdx(idxX, idxY+1)
	lr := g.GetValueAtIdx(idxX+1, idxY+1)

	v:=BilinearInterpolation(&ll,&ul,&lr,&ur,p)

	return float32(v)
}

