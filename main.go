package mygosql

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	f, err := os.Open("data.csv")
	check(err)

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	fmt.Println(line)
	for scanner.Scan() {
		line := scanner.Text()
		parseString(line)
	}

}

func parseString(s string) {
	splited := strings.Split(s, ";")
	fmt.Println(splited)
}
