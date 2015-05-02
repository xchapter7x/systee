package systee

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/arbovm/levenshtein"
)

func DistanceTry() {
	fakeDataSet := []string{
		`2014-02-13T11:44:52.11-0800 [API]     OUT Updated app with guid e1ca6390-cf78-4fc7-9d86-5b7ed01e9c28 ({"instances"=>2})`,
		`2014-02-07T10:54:36.80-0800 [STG]     OUT -----> Downloading and installing node`,
		`2014-02-13T11:44:52.07-0800 [DEA]     OUT Starting app instance (index 1) with guid e1ca6390-cf78-4fc7-9d86-5b7ed01e9c28`,
		`2014-02-13T11:42:31.96-0800 [RTR]     OUT nifty-gui.example.com - [13/02/2014:19:42:31 +0000]
    "GET /favicon.ico HTTP/1.1" 404 23 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_5) AppleWebKit/537.36
    (KHTML, like Gecko) Chrome/32.0.1700.107 Safari/537.36" 10.10.2.142:6609 response_time:0.004092262
    app_id:e1ca6390-cf78-4fc7-9d86-5b7ed01e9c28`,
		`2014-02-13T11:44:27.71-0800 [App/0]   OUT Express server started`,
		"alkjsdf",
		"lkd",
		"aklshdglkahslkdhgas",
		"kkkkkkkasdgkasdg",
		"asdfasdgwaehaweh",
		"sadgweg",
		"dsageeee",
		"agoasdohaheobe",
		"22g2q3g24g",
		"assssssssseeeeeeee",
		"dasgasbeebebebebebe",
		"cvcvcvcvcvcvc",
	}
	km := NewLogKMeans(fakeDataSet, (len(fakeDataSet) / 5))
	km.Group()
	km.Balance()
}

func NewLogKMeans(data []string, centeroids int) (km *LogKMeans) {
	km = &LogKMeans{
		dataSet: data,
	}
	km.setRandomCenteroids(centeroids)
	return
}

type LogKMeans struct {
	dataSet    []string
	centeroids []string
	groups     [][]string
}

func (s *LogKMeans) setRandomCenteroids(setSize int) {
	cnt := 0
	rand.Seed(time.Now().UTC().UnixNano())

	for _, i := range rand.Perm(len(s.dataSet)) {
		fmt.Println(s.dataSet[i])
		s.centeroids = append(s.centeroids, s.dataSet[i])
		s.groups = append(s.groups, []string{})
		cnt++

		if cnt > setSize {
			break
		}
	}
}

func weighted(d, c string) float64 {
	return ((float64(levenshtein.Distance(d, c)) / float64(len(c))) * 10.0)
}

func normal(d, c string) float64 {
	return float64(levenshtein.Distance(d, c))
}

func (s *LogKMeans) Group() {
	var wg sync.WaitGroup

	for _, i := range s.dataSet {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			distance := -1.0
			group := -1

			for i, c := range s.centeroids {

				if ld := normal(d, c); distance == -1 || ld < distance {
					fmt.Println(ld)
					distance = ld
					group = i
				}
			}
			s.groups[group] = append(s.groups[group], d)
		}(i)
	}
	wg.Wait()
	fmt.Println(s.groups, s.centeroids, s.dataSet)
}

func (s *LogKMeans) Balance() {
	fmt.Println("rebalancing...")

	for gi, g := range s.groups {
		sort.Sort(sort.StringSlice(g))
		mid := (len(g) / 2)
		s.centeroids[gi] = g[mid]
		s.groups[gi] = []string{}
	}
	s.Group()
}
