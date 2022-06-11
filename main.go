package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter(wordlist *[]string) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// This route/function takes a URL, shortens it, stores that association, and returns it
	r.POST("/shorten/", func(c *gin.Context) {
		var chosenWords []string
		var numWordsNeeded = 4
		wordlistLength := len(*wordlist)

		// Get post -d "DATA" by grabbing the raw request data
		rawData, err := c.GetRawData()
		if err != nil {
			c.String(http.StatusBadRequest, http.StatusText(400))
			return
		}

		// Parse that; it should just be a valid URI
		originalUrl, err := url.ParseRequestURI(string(rawData))
		if err != nil {
			log.Printf("DEBUG: invalid URL posted: ", err)
			c.String(http.StatusBadRequest, http.StatusText(400))
			return
		}
		log.Printf("DEBUG: url param was ", originalUrl)

		// TODO environment variable for url-wordlength
		for i := 0; i < numWordsNeeded; i++ {
			randIdx := rand.Intn(wordlistLength)
			// NOTE: This dereferencing-the-pointer-before-indexing step confused me
			chosenWords = append(chosenWords, (*wordlist)[randIdx])
		}

		chosenwords_str := strings.Join(chosenWords, "-")
		log.Printf("Chosen words are", chosenwords_str)

		// Store it in our database
		db[chosenwords_str] = originalUrl.String()

		// Return a response
		shortenedUrl := fmt.Sprintf("http://%s/url/%s", c.Request.Host, chosenwords_str)
		c.String(http.StatusOK, shortenedUrl)
	})

	// This route/function looks up a shortened URL
	r.GET("/url/:shortenedUrl", func(c *gin.Context) {
		url := c.Params.ByName("shortenedUrl")

		val, exists := db[url]
		if exists {
			c.String(http.StatusOK, val)
		} else {
			c.String(http.StatusNotFound, "No such URL.")
			return
		}

	})

	return r
}

func main() {
	// Create in-memory wordlist
	var wordlist []string

	// Read wordlist file
	file, err := os.Open("wordlist.txt")
	if err != nil {
		panic("Unable to load wordlist. Exiting.")
	}

	log.Printf("DEBUG: Loading wordlist...")
	startTime := time.Now()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		wordlist = append(wordlist, scanner.Text())
	}
	// Not deferred so we can close it before running the web server
	file.Close()

	elapsed := time.Since(startTime)
	log.Printf("DEBUG: Finished loading wordlist. Time elapsed: ", elapsed)

	if err := scanner.Err(); err != nil {
		log.Printf(err)
	}

	r := setupRouter(&wordlist)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
