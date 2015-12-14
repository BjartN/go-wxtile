package data
import (
	"github.com/bjartn/go-wxtile/tile"
	"fmt"
	"encoding/json"
	"os"
)

func ParseJson(file string) (tile.Grid) {

	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error. Aborting." + err.Error())
		return tile.Grid{}
	}

	gridz := make([] tile.Grid, 1)

	jsonParser := json.NewDecoder(f)
	if err = jsonParser.Decode(&gridz); err != nil {
		fmt.Println("Error. Aborting." + err.Error())
		return   tile.Grid{}
	}

	//fmt.Println("Done")

	if len(gridz) > 0 {
		return gridz[0]
	}

	return tile.Grid{}
}
