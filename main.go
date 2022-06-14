package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_WORDLENGTH    = 4
	DEFAULT_WORDLIST_PATH = "wordlist.txt"
)

// TODO what do I do with these? They're global vars, but not constants.
var dbShortToLong = make(map[string]string)
var dbLongToShort = make(map[string]string)
var wordlist []string
var wordlength = DEFAULT_WORDLENGTH

func setupRouter(wordlength int) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// This route/function takes a URL, shortens it, stores that association, and returns it
	r.POST("/shorten/", func(c *gin.Context) {
		var chosenWords []string
		var numWordsNeeded = wordlength
		wordlistLength := len(wordlist)

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

		// PREVENT DUPLICATION - O(1), but this approach doubles space used for each URL
		// If we're already tracking this URL, just return that URL
		short, exists := dbLongToShort[originalUrl.String()]
		if exists {
			log.Printf("DEBUG: Duplicate URL submitted: %s", originalUrl.String())
			c.String(http.StatusOK, fmt.Sprintf("http://%s/url/%s", c.Request.Host, short))
			return
		}

		// TODO environment variable for url-wordlength
		for i := 0; i < numWordsNeeded; i++ {
			randIdx := rand.Intn(wordlistLength)
			// NOTE: This dereferencing-the-pointer-before-indexing step confused me
			chosenWords = append(chosenWords, (wordlist)[randIdx])
		}

		chosenwords_str := strings.Join(chosenWords, "-")
		log.Printf("Chosen words are", chosenwords_str)

		// Store it in our databases
		dbShortToLong[chosenwords_str] = originalUrl.String()
		dbLongToShort[originalUrl.String()] = chosenwords_str

		// Return a response
		shortenedUrl := fmt.Sprintf("http://%s/url/%s", c.Request.Host, chosenwords_str)
		c.String(http.StatusCreated, shortenedUrl)
	})

	// This route/function looks up a shortened URL
	r.GET("/url/:shortenedUrl", func(c *gin.Context) {
		url := c.Params.ByName("shortenedUrl")

		val, exists := dbShortToLong[url]
		if exists {
			c.Redirect(http.StatusPermanentRedirect, val)
		} else {
			c.String(http.StatusNotFound, "No such URL.")
			return
		}
	})

	return r
}

func main() {

	// Parse environment variables
	wordlistPath, ok := os.LookupEnv("WORDLIST_PATH")
	if !ok {
		log.Printf("Using default WORDLIST_PATH: %v", DEFAULT_WORDLIST_PATH)
		wordlistPath = DEFAULT_WORDLIST_PATH
	}

	// This is a bit weird -- using := further down created a bug (wordlength var scoped to an 'if' block?).
	// so in order to use '=' for a multi-return function (Atoi) I have to declare err here too
	var wordlength int
	var err error
	wordlength_str, _ := os.LookupEnv("WORDLENGTH")
	if wordlength_str == "" {
		log.Printf("Using default wordlength: %v", DEFAULT_WORDLENGTH)
		wordlength = DEFAULT_WORDLENGTH
	} else {
		wordlength, err = strconv.Atoi(wordlength_str)
		if err != nil {
			panic("Invalid WORDLENGTH passed!")
		}
	}

	log.Printf("Wordlength is %v", wordlength)

	// Read wordlist file
	file, err := os.Open(wordlistPath)
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
	log.Printf("DEBUG: Finished loading wordlist. Time elapsed: %v", elapsed)

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error while reading wordlist: %v", err)
	}

	r := setupRouter(wordlength)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
