package googlesearch

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// https://scrapfly.io/blog/how-to-scrape-google/

func Search(q string) {
	queryPart := url.Values{}
	queryPart.Add("q", q)

	apiUrl := "https://www.google.com/search?" + queryPart.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 299 {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("\nGoogleSearch failed: \"%s\"", resp.Status)
	}

	b := new(strings.Builder)
	io.Copy(b, resp.Body)
	log.Print(b.String())
	fmt.Println("success") // todo scrape the result

	err = os.WriteFile("../../data/googlesearch.json", []byte(b.String()), 0666)
	if err != nil {
		panic(err)
	}
}
