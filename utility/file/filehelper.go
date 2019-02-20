package file

import (
	"bufio"
	"container/list"
	"github.com/chzyer/readline"
	"github.com/elastos/Elastos.ELA.Monitor/utility/log"
	"io/ioutil"
	"os"
)

func GetDirectoryFileNames(path string) (fileNames *list.List, err error) {
	dirList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Errorf("read dir error: %+v", err)
		return nil, err
	}

	fileNames = list.New()

	for _, file := range dirList {
		if file.IsDir() == false {
			fileNames.PushBack(file.Name())
		}
	}

	return fileNames, err
}

func ReadFileToLines(path string) (lines *list.List, err error) {
	readLine, err := readline.New(path)
	if err != nil {
		panic(err)
	}

	defer readLine.Close()

	lines = list.New()

	for {
		line, err := readLine.Readline()
		if err != nil { // io.EOF
			break
		}

		//println(line)
		lines.PushBack(line)
	}

	return lines, err
}

func ReadLastLinesFromFile(path string, startLine int64) (lines *list.List, err error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	var lineCount int64 = 1
	lines = list.New()
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		if lineCount > startLine {
			lines.PushBack(fileScanner.Text())
		}
		lineCount++
	}

	defer file.Close()

	return lines, err
}
