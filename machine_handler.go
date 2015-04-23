package systee

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/arbovm/levenshtein"
)

func DistanceTry() {
	km := NewLogKMeans([]string{"alkjsdf",
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
	}, 3)
	km.Group()
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

func (s *LogKMeans) Group() {
	var wg sync.WaitGroup

	for _, i := range s.dataSet {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			distance := -1.0
			group := -1

			for i, c := range s.centeroids {

				if ld := ((float64(levenshtein.Distance(d, c)) / float64(len(c))) * 10.0); distance == -1 || ld < distance {
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

}
