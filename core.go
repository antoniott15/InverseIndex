package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetFiles() []string {
	var files []string
	err := filepath.Walk(DIRECTORY, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}


func scanStopLisT(path string) (map[string]int, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	words := make(map[string]int)
	for scanner.Scan() {
		words[scanner.Text()] = 1
	}

	return words, nil
}



func scanWords(path string, stopWords map[string]int, words map[string][]string) (map[string][]string,map[string]int,error) {
	if path == PATHSTOPLIST {
		return nil,nil ,nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil,nil ,err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	wordsCounter := make(map[string]int)
	for scanner.Scan() {
		if _, ok := stopWords[scanner.Text()]; !ok {
			wrd := strings.ReplaceAll(scanner.Text(),",","")
			words[wrd] = append(words[wrd], path)
			if _, ok := wordsCounter[wrd]; !ok {
				wordsCounter[wrd] = 1
			}else{
				wordsCounter[wrd] = wordsCounter[wrd] + 1
			}
		}
	}

	return words, wordsCounter ,nil
}



type WordList struct {
	name string
	count int
	appearIn []string
}




func WritingResult(result []*WordList) []*WordList {
	var words []*WordList
	for i, value := range result {
		words = append(words, value)
		val := fmt.Sprintf("%s: %d %s",  value.name, value.count, value.appearIn)
		if err := Save(val); err != nil {
			fmt.Println(err)
		}
		if i == 100 {
			return words
		}
	}
	return words
}


func Save(result string) error {
	file, err := os.OpenFile(RESULT, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	dataWrite := bufio.NewWriter(file)

	_, err = io.WriteString(dataWrite, "\n"+result)
	if err != nil {
		return err
	}

	dataWrite.Flush()
	return file.Sync()
}

func initF()  []*WordList {
	files := GetFiles()
	stopList, err := scanStopLisT(PATHSTOPLIST)
	if err != nil {
		panic(err)
	}

	var word []*WordList

	for _, elements := range files {
		resultWords, resultCount ,err := scanWords(elements, stopList, map[string][]string{} )
		if err != nil {
			panic(err)
		}
		for key, value := range resultWords {
			word = append(word, &WordList{
				name:     key,
				count:    resultCount[key],
				appearIn: value,
			})
		}
	}

	sort.Slice(word, func(i, j int) bool {
		return word[i].count > word[j].count
	})


	return WritingResult(word)
}


func QueryMachine(query string, toQuery []*WordList) bool {
	query = strings.ToLower(query)
	str := strings.Split(query, " ")

	var ok, ok2  = false , false
	if str[0] == "not" {
		for _,elements := range toQuery {
			if elements.name == str[1] {
				ok = true
			}
			if elements.name == str[len(str)-1] {
				ok2 = true
			}
		}
	} else {
		for _,elements := range toQuery {
			if elements.name == str[0] {
				ok = true
			}
			if elements.name == str[len(str)-1] {
				ok2 = true
			}
		}
	}

	if !strings.Contains(query, "not") {
		if strings.Contains(query, "or") {
			return ok || ok2
		} else if strings.Contains(query, "and") {

			return ok && ok2
		}
	} else {
		if str[0] == "not" {
			if strings.Contains(query, "or") {
				return !ok || ok2
			} else if strings.Contains(query, "and") {
				return !ok && ok2
			}
		} else {
			if strings.Contains(query, "or") {
				return ok || !ok2
			} else if strings.Contains(query, "and") {
				return ok && !ok2
			}
		}
	}

	return false
}



func main(){
	if err := os.Remove(RESULT); err!= nil {
		panic(err)
	}
	_, err := os.Create(RESULT)
	if err != nil {
		panic(err)
	}
	c := initF()

	query1 := "huye or asd"
	fmt.Print(QueryMachine(query1, c))

	query2 := "jardinero and llegan"
	fmt.Print(QueryMachine(query2, c))

	query3 := "mithril and not funeral"
	fmt.Print(QueryMachine(query3, c))
}
