package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"io/ioutil"
	"strconv"
	"sync"
	"net/http"
)

func main() {
	// Parse command line arguments
	inputFile := flag.String("f", "", "input file containing list of URLs")
	payload := flag.String("p", "", "payload to replace or append parameter values with")
	appendPayload := flag.Bool("a", false, "flag indicating whether to append payload to existing parameter value or replace it entirely")
	mc := flag.String("mc", "", "HTTP response code(s) to filter on (comma-separated)")
	ms := flag.String("ms", "", "String to search for in the response body")
	threads := flag.Int("t", 1, "Number of goroutines to use")
	help := flag.Bool("h", false, "display usage information")
	silent := flag.Bool("s", false, "hide banner")
	flag.Parse()

	if !*silent {
		fmt.Println(`
.__               .__.           
[__) _.._. _.._ _ [__]._  _ | _ ._.
|   (_][  (_][ | )|  |[ )(_]|(/,[  
                         ._|       
		`)
		fmt.Println("\nRemember that bug bounty and security tools should only be used ethically and responsibly.")
		fmt.Println("Misuse of these tools can lead to harm and legal consequences.")
		fmt.Println("Use these tools with caution and obtain permission before performing any testing or analysis.\n")
	}

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: ParamAngler [OPTIONS]\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
        return
    }

	


	codes := make(map[int]bool)
	if *mc != "" {
		for _, code := range strings.Split(*mc, ",") {
			codeInt, err := strconv.Atoi(code)
			if err != nil {
				panic(err)
			}
			codes[codeInt] = true
		}
	}

	jobs := make(chan string)
	results := make(chan string)

	// Read in URLs from either input file or stdin input
	urls := make([]string, 0)
	if *inputFile != "" {
		file, err := os.Open(*inputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				if len(codes) == 0 || codes[resp.StatusCode] {
					if *ms == "" || strings.Contains(string(body), *ms) {
						results <- url
					}
				}
			}
		}()
	}

	// Loop through each URL and replace or append each parameter value separately
	go func() {
		for _, u := range urls {
			parsedUrl, err := url.Parse(u)
			if err != nil {
				panic(err)
			}

			queryValues := parsedUrl.Query()
			for k, v := range queryValues {
				for i := range v {
					tmpValues := make([]string, len(v))
					copy(tmpValues, v)

					if *appendPayload {
						tmpValues[i] += *payload
					} else {
						tmpValues[i] = *payload
					}

					newQueryValues := url.Values{}
					for j, w := range queryValues {
						if j == k {
							newQueryValues[k] = tmpValues
						} else {
							newQueryValues[j] = w
						}
					}

					parsedUrl.RawQuery = newQueryValues.Encode()
					newUrl := parsedUrl.String()
					//fmt.Println(newUrl)
					jobs <- newUrl
				}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for url := range results {
		fmt.Println(url)
	}
}
