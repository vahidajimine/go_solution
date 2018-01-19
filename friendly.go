package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

//Sort by length code
type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}

func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

//Compares the strings and checks is they are friendly
/*The strings are friendly if all the characters contained in the string
Are uniquely bi-directionally mapped,
HOHO and BIBI are friendly as H <-> B and O <-> I
AA and KJ are not because A has 2 mappings
AAA BBBB are no because they also need to be the same length
// */
func isFriendly(string1, string2 string) bool {
	count1 := len(string1)
	count2 := len(string2)
	if count1 != count2 || string1 == "" || string2 == "" {
		return false
	}
	if string1 == string2 {
		return true
	}
	currentMap := map[byte]byte{}
	for i := 0; i < count1; i++ {
		if val, exists := currentMap[string1[i]]; exists {
			if val != string2[i] {
				return false
			}
		} else if val, exists := currentMap[string2[i]]; exists {
			if val != string1[i] {
				return false
			}
		} else {
			currentMap[string1[i]] = string2[i]
			currentMap[string2[i]] = string1[i]
		}
	}
	return true
}

//Returns the count of words with at least 1 friend
func friendlyCountof(wordList []string) int {
	c := 0
	max := len(wordList)
	m := map[int]int{}
	for i := 0; i < max; i++ {
		w := wordList[i]
		if _, exists := m[i]; exists {
			c++
			continue
		}
		for j := 0; j < max; j++ {
			if i == j {
				continue
			}
			if isFriendly(w, wordList[j]) {
				m[i] = j
				m[j] = i
				c++
				break
			}
		}
	}
	return c
}

//Groups strings into slices where all the strings are of same length
func groupWords(wordList []string) [][]string {
	sort.Sort(ByLength(wordList))
	var gs [][]string
	var g []string
	for i := 0; i < len(wordList); i++ {
		g = append(g, wordList[i])
		if i+1 >= len(wordList) {
			gs = append(gs, g)
			g = nil
		} else if len(wordList[i]) < len(wordList[i+1]) {
			gs = append(gs, g)
			g = nil
		}
	}
	return gs
}

//Returns the same result of friendlyCountOf, but uses goroutines and channels to divide the work concurrently
func threadedFriendlyCountOf(wordList []string) int {
	g := groupWords(wordList)
	var c []chan int
	for i := 0; i < len(g); i++ {
		c = append(c, make(chan int))
	}
	for i := 0; i < len(g); i++ {
		go func(wl []string, index int) {
			c[index] <- friendlyCountof(wl)
		}(g[i], i)
	}
	count := 0
	for i := 0; i < len(c); i++ {
		count += <-c[i]
	}
	return count
}

func testIsFriendly() bool {
	fmt.Println("Testing Friendly Func")
	count := 1
	if isFriendly("GAGA", "BOBO") != true {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //2
	if isFriendly("HHHH", "BOBO") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //3
	if isFriendly("HHHH", "HHHH") != true {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //4
	if isFriendly("HHHH", "aeoifnaoienf") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //5
	if isFriendly("JKKJJ", "JKKJ") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //6
	if isFriendly("ABCE", "EFGH") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //7
	if isFriendly("ABCE", "EFGA") != true {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //8
	if isFriendly("", "") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //9
	if isFriendly("", "d") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //10
	if isFriendly("d", "") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	count++ //11
	if isFriendly("EGONUH", "LALALA") != false {
		fmt.Println("Failed Case " + strconv.Itoa(count))
		return false
	}
	fmt.Println("Testing Success")
	return true
}

func testGroupWords() bool {
	//Test case
	p := []string{"LALALA", "XOXOXO", "GCGCGC", "HHHCCC", "BBBMMM", "EGONUH", "HHRGOE", "XOXO", "JUJU", "JKKK", "J", "", ""}
	//Correct grouping
	r := [][]string{[]string{"", ""}, []string{"J"}, []string{"XOXO", "JUJU", "JKKK"}, []string{"LALALA", "XOXOXO", "GCGCGC", "HHHCCC", "BBBMMM", "EGONUH", "HHRGOE"}}
	gp := groupWords(p)
	fmt.Println("Testing groupWords")
	for i := 0; i < len(gp); i++ {
		for j := 0; j < len(gp[i]); j++ {
			//Order doesn't matter
			sort.Strings(gp[i])
			sort.Strings(r[i])
			if gp[i][j] != r[i][j] {
				fmt.Println("Failed on " + gp[i][j] + " " + r[i][j])
				fmt.Println(gp)
				fmt.Println(r)
				return false
			}
		}
	}
	return true
}

func testFriendlyCountOf() bool {
	bob := []string{"LALALA", "XOXOXO", "GCGCGC", "HHHCCC", "BBBMMM", "EGONUH", "HHRGOE", "XOXO", "JUJU", "JKKK", "J", "", ""}
	fmt.Println("Testing FriendlyCount Function")
	v := friendlyCountof(bob)
	if v == 7 {
		return true
	}
	fmt.Println(bob)
	fmt.Print("Failed: Value was ")
	fmt.Print(v)
	fmt.Println("")
	return false
}

func testThreadedFriendlyCountOf() bool {
	p := []string{"LALALA", "XOXOXO", "GCGCGC", "HHHCCC", "BBBMMM", "EGONUH", "HHRGOE", "XOXO", "JUJU", "JKKK", "J", "", ""}
	fmt.Println("Testing Threading")
	v := threadedFriendlyCountOf(p)
	if v == 7 {
		return true
	}
	fmt.Print("Threading: was ")
	fmt.Print(v)
	fmt.Println("")
	return false
}

func main() {
	fmt.Println(testGroupWords())
	fmt.Println(testIsFriendly())
	fmt.Println(testThreadedFriendlyCountOf())
	fmt.Println(testFriendlyCountOf())
	//*
	f, err := ioutil.ReadFile("words.txt")
	if err != nil {
		fmt.Println(err)
	}
	wl := strings.Split(string(f), "\n")
	fmt.Println(threadedFriendlyCountOf(wl))
	// */
}
