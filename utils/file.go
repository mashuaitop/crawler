package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadLine(fileName string) ([]string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return nil,err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				break
			}
			return nil,err
		}
		result = append(result,line)
	}

	return result ,nil
}
