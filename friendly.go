package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

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
		} else {
			currentMap[string1[i]] = string2[i]
			currentMap[string2[i]] = string1[i]
		}
	}
	return true
}

//Returns the friendly number for
func friendlyCountof(wordList []string) int {
	c := 0
	max := len(wordList) - 1
	for i := 0; i < max; i++ {
		w := wordList[0]
		wordList = append(wordList[:0], wordList[1:]...)
		for j := 0; j < len(wordList); j++ {
			if isFriendly(w, wordList[j]) {
				c++
				break
			}
		}
	}
	return c
}

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

func threadedFriendlyCountOf(wordList []string) int {
	g := groupWords(wordList)
	fmt.Println(g)
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

func main() {
	p := []string{"LALALA", "XOXOXO", "GCGCGC", "HHHCCC", "BBBMMM", "EGONUH", "HHRGOE", "XOXO", "JUJU", "JKKK", "J", "", ""}
	fmt.Println(p)
	fmt.Println(isFriendly("GAGA", "BOBO"))    //true
	fmt.Println(isFriendly("HHHH", "BOBO"))    //false
	fmt.Println(isFriendly("HHHH", "HHHH"))    //true
	fmt.Println(isFriendly("HHHH", "oaeinfa")) //false
	fmt.Println(isFriendly("JKKJJ", "JKKJ"))   //false
	fmt.Println(isFriendly("ABCE", "EFGH"))    //false
	fmt.Println(isFriendly("ABCE", "EFGA"))    //true
	fmt.Println(isFriendly("", ""))            //false
	fmt.Println(isFriendly("", "d"))           //false
	fmt.Println(isFriendly("d", ""))           //false
	gp := groupWords(p)
	fmt.Println(gp)
	fmt.Println(threadedFriendlyCountOf(p))
	f, err := ioutil.ReadFile("words.txt")
	if err != nil {
		fmt.Println(err)
	}
	wl := strings.Split(string(f), "\n")
	fmt.Println(len(wl))
	fmt.Println(threadedFriendlyCountOf(wl))
}
