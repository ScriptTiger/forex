package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Function to display help text and exit
func help(err int) {
	os.Stdout.WriteString(
		"Usage: forex [options...]\n"+
		" [-base] <base>         Base currency\n"+
		" [-quote] <quote>       Quote currency\n"+
		" -idecimal <separator>   Input decimal separator\n"+
		" -ithousands <separator> Input thousands seperator\n"+
		" -odecimal <separator>   Output decimal separator\n"+
		" -othousands <separator> Output thousands seperator\n"+
		" -rest <address:port>   Start REST API on given socket\n")
	os.Exit(err)
}

// Precompile regexp pattern
var pattern, _ = regexp.Compile("<div><div class=\"BNeawe iBp4i AP7Wnd\"><div><div class=\"BNeawe iBp4i AP7Wnd\">.*</div></div></div></div></div><div class=\"nXE3Ob\">")

// Query for pair
func query(rest bool, base *string, quote *string, idecimal *string, ithousands *string, odecimal *string, othousands *string) (*string, error) {

	// URL encode white space
	*base = strings.ReplaceAll(strings.ReplaceAll(*base, " ", "+"), "	", "+")
	*quote = strings.ReplaceAll(strings.ReplaceAll(*quote, " ", "+"), "	", "+")

	// Send GET request to query Google for pair
	response, err := http.Get("https://www.google.com/search?q="+*base+"/"+*quote+"&hl=en")
	if err != nil {
		os.Stdout.WriteString("Problem requesting to server")
		help(3)
	}

	// Read query response into string
	responseString, err := io.ReadAll(response.Body)
	if err != nil {
		os.Stdout.WriteString("Problem reading response body")
		help(4)
	}

	// Extract element containing rate from response string
	match := pattern.Find(responseString)
	if match == nil {
		if bytes.Contains(responseString, []byte("Our systems have detected unusual traffic from your computer network.")) {
			return nil, errors.New("The IP address of this forex instance has been flagged by Google for unusual traffic.\nPlease wait for the cool-down period and try again.")
		} else {return nil, errors.New("Pair not found")}
	}

	// Set format defaults if not given

	if idecimal == nil {idecimal = new(string)}
	if ithousands == nil {ithousands = new(string)}
	if odecimal == nil {odecimal = new(string)}
	if othousands == nil {othousands = new(string)}

	if *idecimal == "" {
		if *ithousands == "." {*idecimal = ","
		} else {*idecimal = "."}
	}
	if *ithousands == "" {
		if *idecimal == "," {*ithousands = "."
		} else {*ithousands = ","}
	}
	if *odecimal == "" {*odecimal = *idecimal}

	// Extract rate from element
	rawRate := strings.Split(strings.Split(string(match[76:len(match)-50]), " ")[0], *idecimal)

	// Format rate and return
	rate := strings.ReplaceAll(rawRate[0], *ithousands, *othousands)+*odecimal+rawRate[1]
	return &rate, nil
}

// REST API
func httpHandle(response http.ResponseWriter, request *http.Request) {

	// Get request queries
	base := request.URL.Query().Get("base")
	quote := request.URL.Query().Get("quote")
	idecimal := request.URL.Query().Get("idecimal")
	ithousands := request.URL.Query().Get("ithousands")
	odecimal := request.URL.Query().Get("odecimal")
	othousands := request.URL.Query().Get("othousands")

	// Report error if base or quote not given
	if base == "" || quote == "" {
		response.Write([]byte("Malformed request"))
		return
	}	

	// Query Google for rate
	rate, err := query(true, &base, &quote, &idecimal, &ithousands, &odecimal, &othousands)
	if err != nil {
		response.Write([]byte(err.Error()))
		return
	}

	// Write rate as response
	response.Write([]byte(*rate))
}

func main() {

	// Display help and exit if not enough arguments
	if len(os.Args) < 2 {help(1)}

	// Declare flag pointers
	var (
		base *string
		quote *string
		idecimal *string
		ithousands *string
		odecimal *string
		othousands *string
		rest *string
	)

	// Push arguments to flag pointers
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			switch strings.TrimPrefix(os.Args[i], "-") {
				case "base":
					i++
					if base != nil {help(1)}
					if len(os.Args) == i {help(1)}
					base = &os.Args[i]
					continue
				case "quote":
					i++
					if quote != nil {help(1)}
					if len(os.Args) == i {help(1)}
					quote = &os.Args[i]
					continue
				case "idecimal":
					i++
					if idecimal != nil {help(1)}
					if len(os.Args) == i {help(1)}
					idecimal = &os.Args[i]
					continue
				case "ithousands":
					i++
					if ithousands != nil {help(1)}
					if len(os.Args) == i {help(1)}
					ithousands = &os.Args[i]
					continue
				case "odecimal":
					i++
					if odecimal != nil {help(1)}
					if len(os.Args) == i {help(1)}
					odecimal = &os.Args[i]
					continue
				case "othousands":
					i++
					if othousands != nil {help(1)}
					if len(os.Args) == i {help(1)}
					othousands = &os.Args[i]
					continue
				case "rest":
					i++
					if rest != nil {help(1)}
					if len(os.Args) == i {help(1)}
					rest = &os.Args[i]
					continue
				default:
					help(1)
			}
		} else if base == nil {base = &os.Args[i]
		} else if quote == nil {quote = &os.Args[i]
		} else {help(1)}
	}

	// Initialize REST API if defined
	if rest != nil {
		http.HandleFunc("/api", httpHandle)

		var restAddr string
		if strings.HasPrefix(*rest, ":") && strings.Count(*rest, ":") == 1 {
			restAddr = "localhost"+*rest
		} else {restAddr = *rest}
		os.Stdout.WriteString("Listening on "+restAddr+"\nExample: "+restAddr+"/api?base=eur&quote=usd&idecimal=.&ithousands=,&odecimal=,&othousands=.\n")

		http.ListenAndServe(*rest, nil)
		os.Exit(0)
	}

	// Display help and exit if no base or quote given
	if base == nil || quote == nil {help(2)
	} else {
		rate, err := query(false, base, quote, idecimal, ithousands, odecimal, othousands)
		if err != nil {
			os.Stdout.WriteString(err.Error())
			os.Exit(-1)
		}
		os.Stdout.WriteString(*rate)
		os.Exit(0)
	}
}