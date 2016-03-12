package main


import (

	"os"
	"strings"
	"strconv"
	"github.com/zond/god/common"
	"3_1/Path"
	"fmt"
	"io/ioutil"


)


func getPaths(a int) []Path.Path {
	//b := make([]byte,7000)
	s,_ := os.Getwd()
	dat, _ := ioutil.ReadFile(s + "L" + strconv.Itoa(a) + ".txt")
	var c []Path.Path
	common.MustJSONDecode(dat, &c)
	return c

}

func shortPath(a, b int){

}


func main() {
	path,_ := os.Getwd()

	f, _ := os.Open("C:\\Users\\Mindmaster\\Documents\\GitHub\\go\\src\\3_1\\path\\paths2.txt")
	b := make([]byte,5000)
	_, _ = f.Read(b)
	j := string(b)
	filler := make([][]Path.Path,33)
	
	f.Close()
	s := strings.Split(j, "\n")
	for i := 0; i < len(s); i++ {
		fmt.Println(i)
		ss := strings.Split(s[i], ",")
		n1,_ := strconv.Atoi(ss[0])
		n2,_ := strconv.Atoi(ss[1])
		n3,_ := strconv.Atoi(ss[2])
		p := Path.Path{L1: n1, L2: n2, Length: n3, K: "p" + ss[0] + ss[1] + ".txt"}
		filler[n1] = append(filler[n1],p)
		filler[n2] = append(filler[n2],p)
		/*val := common.MustJSONEncode(p)
		f, _ = os.Create(path + "\\p" + ss[0] + "_" + ss[1] + ".txt")
		f.Write(val)
		f.Close()
		*/

	}

	_, d := Path.ShortestPaths()

	//d = d[]
	//fmt.Println(d[0])
	//fmt.Println(d[1])

	for i := 1; i < len(d); i++ {
		for j := 1; j < len(d[i]); j++ {
			if i != j {
				f, _ = os.Create(path + "\\d" + strconv.Itoa(i) + "_" + strconv.Itoa(j) + ".txt")
				f.Write(common.MustJSONEncode(d[i][j]))
				f.Close()
			}
		}
	}

	/*


	for i := 1; i < 33; i++{
		fmt.Println(i)
		val1, val2, _ := Path.ShortestPath(i, i+1, nil, nil, 0)
		common.MustJSONEncode(Path.Direction{Paths: val1, Dir: val2})
		f, _ = os.Create(path + "\\d" + strconv.Itoa(i) + "_" + strconv.Itoa(i + 1) + ".txt")
	}

	*/

}





/*
import (

	"os"
	"strings"
	"strconv"
	"github.com/zond/god/common"
	"3_1/Path"
	"fmt"


)



func main() {

	path,_ := os.Getwd()

	f, _ := os.Open("C:\\Users\\Mindmaster\\Documents\\GitHub\\go\\src\\3_1\\path\\paths2.txt")
	b := make([]byte,5000)
	_, _ = f.Read(b)
	j := string(b)
	filler := make([][]Path.Path,33)
	
	f.Close()
	s := strings.Split(j, "\n")
	for i := 0; i < len(s); i++ {
		fmt.Println(i)
		ss := strings.Split(s[i], ",")
		n1,_ := strconv.Atoi(ss[0])
		n2,_ := strconv.Atoi(ss[1])
		n3,_ := strconv.Atoi(ss[2])
		p := Path.Path{L1: n1, L2: n2, Length: n3, K: "p" + ss[0] + ss[1] + ".txt"}
		filler[n1] = append(filler[n1],p)
		filler[n2] = append(filler[n2],p)
		val := common.MustJSONEncode(p)
		f, _ = os.Create(path + "\\p" + ss[0] + "_" + ss[1] + ".txt")
		f.Write(val)
		f.Close()

	}

	for i := 1; i < len(filler); i++ {
		loc := Path.Location{Paths: filler[i], N: i, K: "L" + strconv.Itoa(i) + ".txt"}
		f, _ = os.Create(path + "\\" + "L" + strconv.Itoa(i) + ".txt")
		val := common.MustJSONEncode(loc)
		f.Write(val)
		f.Close()
	}


}

*/