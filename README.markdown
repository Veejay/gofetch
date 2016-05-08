# Go Fetch

## Usage
	
Start the server

```bash
./gofetch
```

## Cross-compilation 

Golang allows us to compile the program for a different architecture directly
```
env GOOS=darwin GOARCH=amd64 go build gofetch.go  
```

You can now browse localhost on port 12345

Type in a URL in the input field and click `go`, it will display the results.

## Documentation

Not much to document eh, but well, you can check the [wiki](https://github.com/Veejay/gofetch/wiki)

## What it does

* Gets the URL to search over WebSockets
* Finds all the links on a webpage located at `url` and issues a HTTP GET request on the associated `href`
* Sends back the results as it gets them over the WebSocket back to the client
* In turn the client's onmessage function will add a div containing info about the results inside with 
  appropriate color codes.
* Chokes on some malformed stuff

## What I'd like to add in the near future

* Cache the a MD5 of the already-checked URLs for some time to avoid needless requests (TTL might be nice)
* Take the initial URL from a job queue or something (Redis or whatever)
* Reports (email? JSON?)
