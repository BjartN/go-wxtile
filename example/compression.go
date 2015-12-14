package example

import (
	"encoding/binary"
	"bytes"
	"compress/gzip"
	"github.com/bjartn/go-wxtile/tile"
	"time"
	"fmt"
	"encoding/gob"
	"compress/lzw"
	"compress/flate"
)


func ratio(a int,b int) float64{
	return float64(a) / float64(b)
}

func Compress() {
	fmt.Println("Running...")

	//make grid
	gridSize := 180 * 8 * 360 * 8
	a := make([]float32, gridSize)
	for i := range a {
		a[i] = float32(i)
	}

	//lzw
	compressed := LzwEncode(a)
	fmt.Println("Lzw ratio is ", ratio(len(compressed),(len(a)*4)))

	//un-lzw
	decompressed :=LzwDecode(compressed)

	if decompressed[1000] == a[1000] {
		fmt.Println("LZW is good")
	}

	//gob
	gobbed := GobEncode(a)
	fmt.Println("Gob ratio is ", ratio(len(gobbed),(len(a)*4)))

	//ungob
	ungobbed :=GodDecode(gobbed)

	if ungobbed[1000] == a[1000] {
		fmt.Println("Gob is good")
	}

	//zip
	zipped := GZipEncode(a)
	fmt.Println("Zip ratio is ", ratio(len(zipped),(len(a)*4)))

	//unzip
	unzipped := GZipDecode(zipped)

	if unzipped[1000] == a[1000] {
		fmt.Println("Zip is good")
	}



	//flate
	flated := FlateEncode(a)
	fmt.Println("Deflate ratio is ", ratio(len(flated),(len(a)*4)))

	//deflate
	unflated := FlateDecode(zipped)

	if unflated[1000] == a[1000] {
		fmt.Println("Deflate is good")
	}
}

func GobEncode(a []float32) []byte{
	defer tile.TimeTrack(time.Now(), fmt.Sprint("GobEncode"))

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	encoder.Encode(a)
	return b.Bytes()
}

func GodDecode(a []byte) []float32{
	defer tile.TimeTrack(time.Now(), fmt.Sprint("GodDecoder"))

	b:=bytes.NewBuffer(a)
	decoder := gob.NewDecoder(b)

	var aTarget []float32
	decoder.Decode(&aTarget)

	return aTarget
}

func LzwDecode(a []byte) []float32 {
	defer tile.TimeTrack(time.Now(), fmt.Sprint("LzwDecode"))

	b := bytes.NewBuffer(a)
	r := lzw.NewReader(b,lzw.LSB,8)
	defer r.Close()

	aTarget:=make([]float32, len(a)/4)

	for i := range aTarget {

		var f float32
		binary.Read(r, binary.LittleEndian,&f)
		aTarget[i]=f
	}

	return aTarget
}


func LzwEncode(a []float32) []byte{
	defer tile.TimeTrack(time.Now(), fmt.Sprint("LzwEncode"))

	var b bytes.Buffer
	w := lzw.NewWriter(&b,lzw.LSB,8)
	defer w.Close()

	for i := range a {
		binary.Write(w, binary.LittleEndian, a[i])
	}

	return b.Bytes()
}


func GZipEncode(a []float32) []byte{
	defer tile.TimeTrack(time.Now(), fmt.Sprint("GZipEncode"))

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	for i := range a {
		binary.Write(w, binary.LittleEndian, a[i])
	}

	return b.Bytes()
}


func GZipDecode(a []byte) []float32 {
	defer tile.TimeTrack(time.Now(), fmt.Sprint("GZipDecode"))

	b := bytes.NewBuffer(a)
	r,_ := gzip.NewReader(b)
	defer r.Close()

	aTarget:=make([]float32, len(a)/4)

	for i := range aTarget {

		var f float32
	   	binary.Read(r, binary.LittleEndian,&f)
		aTarget[i]=f
	}

	return aTarget
}


func FlateEncode(a []float32) []byte{
	defer tile.TimeTrack(time.Now(), fmt.Sprint("FlateEncode"))

	var b bytes.Buffer
	w,_ := flate.NewWriter(&b,1)
	defer w.Close()

	for i := range a {
		binary.Write(w, binary.LittleEndian, a[i])
	}

	return b.Bytes()
}


func FlateDecode(a []byte) []float32 {
	defer tile.TimeTrack(time.Now(), fmt.Sprint("FlateDecode"))

	b := bytes.NewBuffer(a)
	r := flate.NewReader(b)
	defer r.Close()

	aTarget:=make([]float32, len(a)/4)

	for i := range aTarget {

		var f float32
		binary.Read(r, binary.LittleEndian,&f)
		aTarget[i]=f
	}

	return aTarget
}


