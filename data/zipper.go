package data
import (
	"fmt"
	"archive/zip"
)

//ListZip prints the contents of a zip file
func ListZip(file string) {

	r, err := zip.OpenReader(file)
	if (err != nil) {
		fmt.Println("Error. Aborting." + err.Error())
		return
	}
	defer r.Close() //call after this function closes

	for _, f := range r.File {

		if (f.FileInfo().IsDir()) {
			continue
		}

		fmt.Println(f.Name)


	}

	fmt.Println("Zip done")
}
