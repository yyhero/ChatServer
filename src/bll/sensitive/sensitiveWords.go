package sensitive

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	sensitiveWordsList = make([]string, 0, 2048)
)

func init() {
	if tmpList, err := readSensitiveWords(); err != nil {
		panic(err)
	} else {
		sensitiveWordsList = tmpList
	}
}

func readSensitiveWords()  ( []string, error){
	fileName := "sensitive.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil, err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	list := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		list = append(list, line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil ,err
			}
		}
	}
	return list, nil
}

func HandleSensitiveWords(input string) string {
	if len(sensitiveWordsList) == 0 {
		return input
	}
	for _, item := range sensitiveWordsList {
		if strings.Contains(strings.ToLower(input), item) {
			input = strings.Replace(input, item, "***", -1)
		}
	}
	return input
}


