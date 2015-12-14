package tile
import (
	"time"
	"log"
)

//TimeTrack is a utility function used for timing the execution of a function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}