# Golang URL Shortener

Just a fun way to practice Go.


## Setup Instructions
```
cd $PROJECT_ROOT
cat /usr/share/dict/words > wordlist.txt
go run main.go
```

## Usage

### Shorten a URL
```
curl -X POST -d "https://tutorialinux.com" localhost:8080/shorten
```

### Resolve a shortened URL
Paste the URL that was returned from your original curl request into a browser. Or:
```
curl http://localhost:8080/url/silverberry-foppy-betocsin-underfrock
```


## TODO
- duplicate detection for generated short-urls (re-roll mechanic)
  - shorter default! (2 or 3 words?)
    - Should be ok -- 235886 words in wordlist, ^2 == 55,642,204,996 URLs before we start having problems
- finish writing some tests
- create a profile from those tests
- in the "/shorten/" handler, extract all the request-data-into-URL munging into its own function
- database: optionally make persistent? Use Redis? Maybe like ./myapp --redis=redis.local:6379 cues the application into the fact that we want to use redis for persistent storage, and not our little baby in-memory map DB.
- environment var: Log level (DEBUG, ERR). Maybe [something fancy like zap](https://github.com/uber-go/zap)?
  - Existing loglines should be categorized into those log levels


### Deduping or Dupe-Avoidance:

- *[CURRENT IMPLEMENTATION]* we could avoid dupes FAST (at the expense of space) with a second map that does the reverse association ({origURL:shortenedURL}). Every storage operation could then check to see if there's already a shortened URL cached. Cost: 2n space, not terrible, maybe kinda bad in practice.

- HOWEVER we could cheat: if we only store {origURL:shortenedURL} (which would cause lookups to be slow - O(n)), then we could cache those slow lookups aggressively via the webserver/CDN fronting the app and not worry about it. A cache warming step in a separate goroutine could be a super hilarious/gross hacky way of doing this. It's sick, really, but it could work as a cheat to get Theta(1)-like time without using twice the space.

- actual dedupes in an async goroutine to periodically "clean up" the main datastructure. O(n), but in a controlled way, and without doubling the space needed.
