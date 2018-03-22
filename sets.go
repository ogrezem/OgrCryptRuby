package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
	"os"
	"math"
)

const (
	isDebug bool = false
)

var (
	notations = map[int]bool {
		2 : true,
		10 : true,
	}
)

func main(){
	start:
	for {
		var num int
		var stringNum string
		fmt.Print("Введите число входных переменных: ")
		_, err := fmt.Scan(&stringNum)
		if err != nil {
			panic(err)
		}
		isInteger, err := regexp.MatchString(`^[0-9]*$`, stringNum)
		if err != nil {
			panic(err)
		} else if !isInteger {
			fmt.Println("Число входных переменных может принимать лишь целочисленное значение.")
			continue
		}
		num, _ = strconv.Atoi(stringNum)
		var setsAmount int = NaturStepen(2, num)
		sets := make(map[string]int)
		for i := 0; i < setsAmount; i++ {
			sets[IntToBinaryWithSize(num, i)] = 0
		}
		if isDebug {
			fmt.Println(sets)
		}
		var notation int
		for {
			fmt.Print("В какой системе счисления собираетесь выбирать наборы? ")
			fmt.Scan(&notation)
			if !notations[notation] {
				fmt.Println("Вы не можете выбрирать наборы в этой системе счисления. Попробуйте ещё раз.")
				continue
			} else {
				break
			}
		}
		for {
			var selectedSet string
			if notation == 10 {
				fmt.Print("Напишите порядковый номер набора, который хотите просмотреть, либо введите exit, если вы уже определились со значениями на каждом из наборов: ")
			} else if notation == 2 {
				fmt.Print("Введите набор в двоичном коде, который хотите посмотреть, либо введите exit, если вы уже определились со значениями на каждом из наборов: ")
			}
			fmt.Scan(&selectedSet)
			var intSelectedSet int
			if selectedSet == "exit" {
				break
			} else {
				hasLiterals, err := regexp.MatchString(`[a-z]|[A-Z]|\.|\,+`, selectedSet)
				if err != nil {
					panic(err)
				} else if hasLiterals {
					fmt.Println("Ошибка. Можно вводить только числа и команды!")
					continue
				} else {
					if notation == 2 {
						isBinary, err := regexp.MatchString(`^[01]*$`, selectedSet)
						if err != nil {
							panic(err)
						} else if !isBinary {
							fmt.Println("Введённая строка не может являться набором.")
							continue
						}
					}
					fmt.Println("Выбранный набор: " + selectedSet)
					if notation == 10 {
						intSelectedSet, err = strconv.Atoi(selectedSet)
						if err != nil {
							panic(err)
						} else if intSelectedSet >= setsAmount {
							fmt.Println("Такого набора нет в заданной коллекции наборов")
							continue
						}
					} else if notation == 2 {
						if len(selectedSet) > num {
							fmt.Println("Такого набора нет в заданной коллекции наборов")
							continue
						} else if len(selectedSet) < num {
							selectedSet = ConvertToSize(num, selectedSet)
						}
					}
				}
			}
			var binarySeenSet string
			if notation == 10 {
				binarySeenSet = IntToBinaryWithSize(num, intSelectedSet)
			} else if notation == 2 {
				binarySeenSet = selectedSet
			}
			fmt.Print("Набор: ", binarySeenSet, ". Задайте значение функции при данном наборе, введя 0 или 1: ")
			var answer string
			fmt.Scan(&answer)
			fmt.Println("Ваш ответ:", answer)
			if answer == "0" || answer == "1" {
				sets[binarySeenSet], _ = strconv.Atoi(answer)
				fmt.Println(binarySeenSet, sets[binarySeenSet])
			} else {
				fmt.Println("Неверная команда")
			}
		}
		sortedSets := make([]string, setsAmount)
		sortedFuncValues := make([]int, setsAmount)
		for set, valueOfSet := range sets {
			decimalSet, err := strconv.ParseInt(set, 2, 64)
			if err != nil {
				panic(err)
			}
			sortedSets[decimalSet] = set
			sortedFuncValues[decimalSet] = valueOfSet
		}
		var choise int
		for {
			fmt.Print("Выберите способ вывода наборов:\n" + "1. На экран\n" + "2. В файл\n" + "3. На экран и в файл\n" +
				"Выбор: ")
			_, err = fmt.Scan(&choise)
			if err != nil {
				panic(err)
			}
			switch choise {
			case 1:
				for i := 0; i < setsAmount; i++ {
					fmt.Println(sortedSets[i], sortedFuncValues[i])
				}
				goto finalChoise
			case 2:
				OutputSetsToFile(sortedSets, sortedFuncValues)
				goto finalChoise
			case 3:
				for i := 0; i < setsAmount; i++ {
					fmt.Println(sortedSets[i], sortedFuncValues[i])
				}
				OutputSetsToFile(sortedSets, sortedFuncValues)
				goto finalChoise
			default:
				fmt.Println("Неверно выбран вариант ответа.")
				continue
			}
		}
		finalChoise:
		var answer string
		fmt.Print("Вы хотите составить ещё таблицу? ")
		for {
			fmt.Scan(&answer)
			switch answer {
			case "да":
				goto start
			case "нет":
				goto exit
			default:
				fmt.Println("Неверно введена команда. Попробуйте ещё раз: ")
				continue
			}
		}
	}
	exit:
}

/*
	Реазиловать шифрование
*/

func OutputSetsToFile(sortedSets []string, sortedFuncValues []int) {
	setsAmount := len(sortedSets)
	if os.Remove("set.txt") != nil {
		fmt.Println("Закройте файл set.txt и повторите попытку")
		os.Exit(0)
	}
	file, err := os.OpenFile("set.txt", os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("Наборы при %d переменных:\n", int(math.Log2( float64( setsAmount )))))
	for i := 0; i < setsAmount; i++ {
		file.WriteString(fmt.Sprintf("%s %d\n", sortedSets[i], sortedFuncValues[i]))
	}
}

func ConvertToSize(neededSize int, lowSizeString string) string {
	if len(lowSizeString) < neededSize {
		var difference int = neededSize - len(lowSizeString)
		return strings.Repeat("0", difference) + lowSizeString
	}
	return lowSizeString
}

func IntToBinaryWithSize(neededSize int, number int) string {
	return ConvertToSize(neededSize, strconv.FormatInt(int64(number), 2))
}

func PrintTable(num int) {
	var iters int = NaturStepen(2, num)
	for i := 0; i < int(iters); i++ {
		fmt.Printf("%" + "00" + strconv.Itoa(num) + "b\n", i)
	}
}

func MakeTable(num int) {
	var iters int = NaturStepen(2, num)
	//numSet := []string{}
	for i := 0; i < int(iters); i++ {
		fmt.Printf("%" + "00" + strconv.Itoa(num) + "b\n", i)
	}
}

func NaturStepen(a int, b int) int {
	var result int = 1
	for i := 0; i < b; i++ {
		result *= a
	}
	return result
}