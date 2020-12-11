package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurlx/pkg/sigurlx"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	threads int
	output  string
	silent  bool
	noColor bool
	verbose bool
}

var (
	co options
	au aurora.Aurora
	so sigurlx.Options
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
     _                  _      
 ___(_) __ _ _   _ _ __| |_  __
/ __| |/ _`+"`"+` | | | | '__| \ \/ /
\__ \ | (_| | |_| | |  | |>  < 
|___/_|\__, |\__,_|_|  |_/_/\_\ v1.0.0
       |___/
`).Bold())
}

func init() {
	// general options
	flag.BoolVar(&co.noColor, "nc", false, "")
	flag.BoolVar(&co.silent, "s", false, "")
	flag.IntVar(&co.threads, "t", 50, "")
	flag.BoolVar(&co.verbose, "v", false, "")

	// task options
	flag.BoolVar(&so.Categorize, "cat", false, "")
	flag.BoolVar(&so.ScanParam, "param-scan", false, "")
	flag.BoolVar(&so.Request, "request", false, "")

	// Http options
	flag.StringVar(&so.UserAgent, "UA", "", "")

	// OUTPUT
	flag.StringVar(&co.output, "o", "", "")

	flag.Usage = func() {
		banner()

		h := "Usage:\n"
		h += "  sigurlx [OPTIONS]\n\n"

		h += "TASK OPTIONS:\n"
		h += "   -cat                       categorize urls\n"
		h += "   -param-scan                scan url parameters\n"
		h += "   -request                   send HTTP request\n\n"

		h += "GENERAL OPTIONS:\n"
		h += "   -t                         number of concurrent threads. (default: 50)\n"
		h += "   -nc                        no color mode\n"
		h += "   -s                         silent mode\n"
		h += "   -v                         verbose mode\n\n"

		h += "REQUEST OPTIONS (used with -request):\n"
		h += "   -UA                        HTTP user agent\n\n"

		h += "OUTPUT OPTIONS:\n"
		h += "   -o                         output file\n\n"

		fmt.Fprintf(os.Stderr, h)
	}

	flag.Parse()

	au = aurora.NewAurora(!co.noColor)
}

func main() {
	if !gos.HasStdin() {
		os.Exit(1)
	}

	if !co.silent {
		banner()
	}

	options, err := sigurlx.ParseOptions(&so)
	if err != nil {
		log.Fatalln(err)
	}

	URLs := make(chan string, co.threads)

	go func() {
		defer close(URLs)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			URLs <- scanner.Text()
		}
	}()

	var output []sigurlx.Results

	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 0; i < co.threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			runner, err := sigurlx.New(options)
			if err != nil {
				log.Fatalln(err)
			}

			for URL := range URLs {
				if URL == "" {
					continue
				}

				results, err := runner.Process(URL)
				if err != nil {
					if co.verbose {
						fmt.Println(err)
					}

					continue
				}

				fmt.Println(results.URL)

				mutex.Lock()
				output = append(output, results)
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()

	if co.output != "" {
		if err := saveResults(co.output, output); err != nil {
			log.Fatalln(err)
		}
	}
}

func saveResults(outputPath string, output []sigurlx.Results) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		directory, _ := path.Split(outputPath)

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			if directory != "" {
				err = os.MkdirAll(directory, os.ModePerm)
				if err != nil {
					return err
				}
			}
		}
	}

	outputJSON, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return err
	}

	outputFile, err := os.Create(co.output)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	_, err = outputFile.WriteString(string(outputJSON))
	if err != nil {
		return err
	}

	return nil
}
