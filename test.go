package main

import (
	"cluster/lib/querymap"
	"cluster/lib/querymap/ops"
	"log"
)

func main() {

	que := querymap.New()
	que.SubAssign("wheels", "rims", 2)
	que.SubAssign("wheels1", "rims", 3)
	que.SubAssign("wheels2", "rims", 7)
	que.SubAssign("wheels3", "rims", 4)
	que.SubAssign("wheels4", "rims", 9)
	que.SubAssign("wheels5", "rims", 34)
	que.SubAssign("wheels6", "rims", -81)
	temp, _ := que.Mask("rims", ops.GT, 3)
	for key, row := range temp.Items() {
		log.Print(key + ": ")
		for subkey, data := range row {
			log.Print(subkey, ": ", data, " ")
		}
	}
	log.Println()
	for key, row := range que.Items() {
		log.Print(key + ": ")
		for subkey, data := range row {
			log.Print(subkey, ": ", data, " ")
		}
	}
	//for {
	//	timel.RemoveBefore(time.Now().UTC().Add(-5000 * time.Microsecond))
	//	time.Sleep(1000 * time.Microsecond)
	//	log.Println(timel.Count(), timel.Start().Equal(timel.First().Time()))
	//}

}
