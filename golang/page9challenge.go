package main

import (
    "crypto/tls"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "regexp"
    "strconv"
)


func quickGuess(body string) (url string, status bool) {
    status = false
    //re, _ := regexp.Compile(`src=\"(/incoming/.*IMAGE_ALTERNATES/p)(\d+)([\w/-]+)\"XX`)
    re, _ := regexp.Compile(`src=\"(/incoming/.*IMAGE_ALTERNATES/p)(\d+)([\w/-]+)\"`)

    match := re.FindStringSubmatch(string(body))
    if match != nil {
        guess := "https://ekstrabladet.dk" + match[1] + "1600" + match[3]

        head, _ := http.Head(guess)
        if head.StatusCode == http.StatusOK {
            url = guess
            status = true
        }
    }
    if verbose ==true {
        if url != "" {
            fmt.Println("DEBUG: function quickGuess() found: " + url)
        } else {
            fmt.Println("DEBUG: function quickGuess() had no luck.")
        }
    }
    return url, status
}


func findMaxRes(curMaxUrl string) (newMaxUrl string, findMaxResStatus bool) {
    findMaxResStatus = false
    newMaxUrl = curMaxUrl

    re, _ := regexp.Compile(`([\w-/]+IMAGE_ALTERNATES/p)+(\d+)(/[\w/-]+)`)

    if re.MatchString(curMaxUrl) {
        gurl_split := re.FindStringSubmatch(curMaxUrl)
        size := gurl_split[2]
        curiter, _ := strconv.ParseFloat(size, 64)

        for {
            curiter += 20
            u := "https://ekstrabladet.dk" + gurl_split[1] + fmt.Sprintf("%v", curiter) + gurl_split[3]
            response, _ := http.Head(u)

            if response.StatusCode == http.StatusOK {
                findMaxResStatus = true
                newMaxUrl = u
            }

            // let's not increase forever
            if curiter >= 4000 {
                break
            }
        }
    }

    if verbose ==true {
        fmt.Println("DEBUG: function findMaxRes() found: " + newMaxUrl)
    }
    return newMaxUrl, findMaxResStatus
}

func findMaxBytes(body string) (url string, status bool) {
    status = false
    re, _ := regexp.Compile(`src=\"([\w/-]+)\"`)

    gurls := re.FindAllStringSubmatch(string(body), -1)

    maxFoundBytes := 0
    maxFoundUrl := ""

    for _, girl := range gurls {
        u := ("https://ekstrabladet.dk" + girl[1])
        response, _ := http.Head(u)

        if response.StatusCode == http.StatusOK {
            status = true
            imglen, _ := strconv.Atoi(response.Header["Content-Length"][0])
            if imglen > maxFoundBytes {
                maxFoundBytes = imglen
                maxFoundUrl = girl[1]
            }
        }
    }
    if verbose ==true {
        fmt.Println("DEBUG: function findMaxBytes() found: " + maxFoundUrl)
    }
    return maxFoundUrl, status
}


var verbose bool

func main() {
    flag.BoolVar(&verbose, "debug", false, "execute with verbosity")
    flag.Parse()

    // Set a fall-back url:
    url := "https://www.youtube.com/watch?v=IO9XlQrEt2Y"

    // Ignore insecure certificates
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}

    req, err := http.NewRequest("GET", "https://ekstrabladet.dk/side9/", nil)
    if err != nil {

        fmt.Println(url)
        os.Exit(1)
    }
    //req.Header.Set("Host", "ekstrabladet.dk")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
    resp, err := client.Do(req)


    // Get the body html:
    //resp, err := client.Get("https://ekstrabladet.dk/side9/")
    if err != nil {
        fmt.Println(url)
        os.Exit(1)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        if verbose ==true {
            fmt.Println("Cannot parse body.")
        }
        fmt.Println(url)
        os.Exit(1)
    }


    // Check if we can do a quick guess:
    quickGuessUrl, quickGuessStatus := quickGuess(string(body))
    if quickGuessStatus ==true {
        url = quickGuessUrl

    } else {
        // Find the largest image size referenced in the body:
        findMaxBytesUrl, findMaxBytesStatus := findMaxBytes(string(body))
        if findMaxBytesStatus ==true {
            url = findMaxBytesUrl
        }
        // Try to guess an even larger img url:
        findMaxResUrl, findMaxResStatus := findMaxRes(findMaxBytesUrl)
        if findMaxResStatus ==true {
            url = findMaxResUrl
        }
    }
    fmt.Println(url)
}
