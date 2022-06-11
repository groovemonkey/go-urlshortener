# Golang URL Shortener

Just a fun way to practice Go.

## Setup Instructions
cd $PROJECT_ROOT
cat /usr/share/dict/words > wordlist.txt
go run main.go


## Usage

```
# Shorten a URL
curl -X POST -d "https://tutorialinux.com" localhost:8080/shorten

# Resolve a shortened URL
curl http://localhost:8080/url/silverberry-foppy-betocsin-underfrock
```

## TODO
- actually, make /url/ do a redirect
- environment var: pass in the wordlist path
- environment var: shortened URL wordlength (right now it's 5)
- database: optionally make persistent? Use Redis? Maybe like ./myapp --redis=redis.local:6379 cues the application into the fact that we want to use redis for persistent storage, and not our little baby in-memory map DB.
- environment var: Log level (DEBUG, ERR). Maybe [something fancy like zap](https://github.com/uber-go/zap)?
  - Existing loglines should be categorized into those log levels

### Deduping or Dupe-Avoidance:

  - we could avoid dupes FAST (at the expense of space) with a second map that does the reverse association ({origURL:shortenedURL}). Every storage operation could then check to see if there's already a shortened URL cached. Cost: 2n space, not terrible, maybe kinda bad in practice.

  - HOWEVER we could cheat: if we only store {origURL:shortenedURL} (which would cause lookups to be slow - O(n)), then we could cache those slow lookups aggressively via the webserver/CDN fronting the app and not worry about it. A cache warming step in a separate goroutine could be a super hilarious/gross hacky way of doing this. It's sick, really, but it could work as a cheat to get Theta(1)-like time without using twice the space.
