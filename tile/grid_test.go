package tile

import (
	"testing"
)

func SetUp() Grid {

	d := []float32{
		8, 1, 1, 1, 1, 1,
		2, 7, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4,
		5, 5, 5, 5, 9, 5,
		6, 6, 6, 6, 6, 12}

	h := GridHeader{
		Nx:       6,
		Ny:       6,
		Dx:       200,
		Dy:       400,
		ScanMode: 0,
		Lo1:      0,
		La1:      2000,
		Lo2:      1000,
		La2:      0}

	return Grid{d, h}

}

func TestShouldScaleUpsideDown(t *testing.T) {

	scale := Scale(0, 10, 100, 0)

	AssertAreEqual(t, 100, scale(0), "Scale upside down")

}

func TestShouldScale(t *testing.T) {

	scale := Scale(0, 10, 0, 100)

	AssertAreEqual(t, 50, scale(5), "Scale mapping from 0-10 to 0-100, should give 50 when scaling 5")

}

func TestShouldFindCoordInGridByIndex(t *testing.T) {
	g := SetUp()

	v := g.GetValueAtIdx(0, 0)
	AssertAreEqual(t, 0, float64(v.X), "X Coord 1")

	v2 := g.GetValueAtIdx(1, 0)
	AssertAreEqual(t, 200, float64(v2.X), "X Coord 2")

	v3 := g.GetValueAtIdx(0, 0)
	AssertAreEqual(t, 2000, float64(v3.Y), "Y Coord 1")

	v4 := g.GetValueAtIdx(0, 1)
	AssertAreEqual(t, 1600, float64(v4.Y), "Y Coord 2")

}

func TestShouldFindCoordWhenAccessingOutsideGridByIndex(t *testing.T) {

	g := SetUp()
	v := g.GetValueAtIdx(5, 6)

	AssertAreEqual(t, -400, float64(v.Y), "Correct coor outside grid")
}

func TestShouldFindLowerLeftValueInGridByCoords(t *testing.T) {

	g := SetUp()

	v2 := g.GetValueAt(&Point{1000, 0})
	AssertAreEqual(t, 12, float64(v2), "Should return lower left value in grid when askin by point")

}

func TestShouldFindUpperRightValueInGridByCoords(t *testing.T) {

	g := SetUp()

	v1 := g.GetValueAt(&Point{0, 2000})
	AssertAreEqual(t, 8, float64(v1), "Should return upper right value in grid when askin by point")

}

func TestShouldFindWidthOfGrid(t *testing.T) {
	g := SetUp()
	AssertAreEqual(t, 1000, g.Width(), "Should return width of grid")
}

func TestShouldFindHeightOfGrid(t *testing.T) {
	g := SetUp()
	AssertAreEqual(t, 2000, g.Height(), "Should return width of grid")
}

func TestShouldFindValueInGridByIndex(t *testing.T) {
	g := SetUp()

	v := g.GetValueAtIdx(1, 1)
	AssertAreEqual(t, 7, float64(v.Value), "Should return value in grid when askin by index")

	v4 := g.GetValueAtIdx(4, 4)
	AssertAreEqual(t, 9, float64(v4.Value), "Should return value in grid when askin by index")

}
