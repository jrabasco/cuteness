package files

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func ListDirectories(folder string) []string {
	res := []string{}
	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return res
	}

	for _, f := range files {
		if f.IsDir() {
			res = append(res, f.Name())
		}
	}

	return res
}

func RandomFile(folder string) string {
	files, err := ioutil.ReadDir(folder)

	if err != nil {
		return ""
	}

	poss := []string{}

	for _, f := range files {
		if !f.IsDir() {
			poss = append(poss, f.Name())
		}
	}

	rand.Seed(int64(time.Now().Nanosecond()))

	return poss[rand.Intn(len(poss))]
}

func RandomFileContents(folder string) string {
	file := RandomFile(folder)

	filePath := fmt.Sprintf("%s/%s", folder, file)
	contents, err := ioutil.ReadFile(filePath)

	if err != nil {
		return ""
	}

	return string(contents)
}
