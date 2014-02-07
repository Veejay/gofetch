# Go Fetch

## Usage
	
It's a simple command, no rocket science involved:
	
	gofetch url

## What it does

* Finds all the links on a webpage located at `url` and issues a HTTP GET request on the associated `href`
* Displays the links and the HTTP status code returned from the GET request
* Chokes on some malformed stuff

## What I'd like to add in the near future

* Cache the a MD5 of the already-checked URLs for some time to avoid needless requests (TTL might be nice)
* Take the initial URL from a job queue or something (Redis or whatever)
* HTTP interface (see gorilla/mux or net/http)
* Reports (email? JSON?)