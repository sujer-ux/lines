package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

var (
	quietMode     = flag.Bool("q", false, "Тихий режим (только итог)")
	skipEmpty     = flag.Bool("ss", false, "Пропускать пустые строки")
	recursive     = flag.Bool("r", false, "Рекурсивный поиск")
	excludeDirs   = flag.String("exclude", "vendor", "Директории для исключения (через запятую)")
)

func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		line := scanner.Text()
		if *skipEmpty {
			if strings.TrimFunc(line, unicode.IsSpace) == "" {
				continue
			}
		}
		lines++
	}
	return lines, scanner.Err()
}

func isExcluded(path string) bool {
	excluded := strings.Split(*excludeDirs, ",")
	for _, dir := range excluded {
		dir = strings.TrimSpace(dir)
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Укажите шаблон (например, *.go)")
		os.Exit(1)
	}

	pattern := args[0]
	var files []string

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil || isExcluded(path) {
			return nil
		}
		if !info.IsDir() {
			matched, _ := filepath.Match(pattern, filepath.Base(path))
			if matched {
				files = append(files, path)
			}
		}
		return nil
	}

	if *recursive {
		filepath.Walk(".", walker)
	} else {
		entries, _ := os.ReadDir(".")
		for _, entry := range entries {
			if !entry.IsDir() && !isExcluded(entry.Name()) {
				matched, _ := filepath.Match(pattern, entry.Name())
				if matched {
					files = append(files, entry.Name())
				}
			}
		}
	}

	if len(files) == 0 {
		if !*quietMode {
			color.Red("Файлы не найдены.")
		}
		return
	}

	total := 0
	blue := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, file := range files {
		lines, err := countLines(file)
		if err != nil {
			if !*quietMode {
				fmt.Printf("%s: %s\n", blue(file), red("ошибка"))
			}
			continue
		}
		if !*quietMode {
			fmt.Printf("%s: %d\n", blue(file), lines)
		}
		total += lines
	}

	if *quietMode {
		fmt.Println(total)
	} else {
		color.Green("Всего строк: %d", total)
	}
}