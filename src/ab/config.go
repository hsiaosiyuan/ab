package ab

import (
	"flag"
	"net/url"
	"log"
	"fmt"
	"strings"
	"errors"
	"os"
)

type Headers [][]string

func (h *Headers) String() string {
	return fmt.Sprintf("%d", h)
}

func (h *Headers) Set(v string) error {
	if kv := strings.Split(v, ":"); len(kv) != 2 {
		return errors.New("Invalid header: " + v)
	}else {
		*h = append(*h, kv)
		return nil
	}
}

type Cookies [][]string

func (c *Cookies) String() string {
	return fmt.Sprintf("%d", c)
}

func (c *Cookies) Set(v string) error {
	if kv := strings.Split(v, "="); len(kv) != 2 {
		return errors.New("Invalid cookie: " + v)
	}else {
		*c = append(*c, kv)
		return nil
	}
}

type Post string

func (p *Post) String() string {
	return string(*p)
}

func (p *Post) Set(v string) error {
	if f, err := os.OpenFile(v, os.O_RDONLY, os.FileMode(0666)); err != nil {
		return err
	}else {
		defer f.Close()
		*p = Post(v)
		return nil
	}
}

type Config struct {
	Number        int     // number of requests to perform for the benchmarking session
	Concurrence   int     // number of multiple requests to perform at a time
	Endpoint      string  // url of the endpoint for testing
	EscapeCache   string  // parameter name will be appended in query string to escape cache. default '_ec=timestamp'
	Cookies       Cookies // cookie. add cookie to the request
	Headers       Headers // custom header, append extra headers to the request
	Post          Post    // file containing data to POST, remember to also set -T
	ContentType   string  // content-type header to use for POST/PUT data, eg. application/x-www-form-urlencoded. default: application/x-www-form-urlencoded.
	TimeLimit     int     // time limit. maximum number of seconds to spend for benchmarking. default: 3600

	EndPointParsed *url.URL
}

var Cfg *Config

func ParseConfig() {
	Cfg = new(Config)

	flag.IntVar(&Cfg.Number, "n", 1, "Number of requests to perform for the benchmarking session")
	flag.IntVar(&Cfg.Concurrence, "c", 1, "Number of multiple requests to perform at a time")
	flag.IntVar(&Cfg.TimeLimit, "t", 3600, "Time limit. maximum number of seconds to spend for benchmarking")

	flag.StringVar(&Cfg.Endpoint, "E", "", "Url of the endpoint for testing")
	flag.StringVar(&Cfg.EscapeCache, "e", "", "Parameter name will be appended in query string to escape cache. eg: '_ec=timestamp'")

	flag.Var(&Cfg.Cookies, "C", "Cookie. add cookie to the request (repeatable)")
	flag.Var(&Cfg.Headers, "H", "Custom header, append extra headers to the request (repeatable)")
	flag.Var(&Cfg.Post, "p", "File containing data to POST, remember to also set -T")
	flag.StringVar(&Cfg.ContentType, "T", "application/x-www-form-urlencoded", "Content-type header to use for POST/PUT data, eg. application/x-www-form-urlencoded. default: text/plain.")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		log.Fatal("None parameter!")
	}

	if Cfg.Endpoint == "" {
		flag.Usage()
		log.Fatal("Provide a endpoint address to test!")
	}

	if _url , err := url.Parse(Cfg.Endpoint); err != nil {
		log.Fatal("Invalid url to test: ", Cfg.Endpoint)
	}else {
		Cfg.EndPointParsed = _url

		if Cfg.EscapeCache != "" {
			fmt.Println("Endpoint to test:", _url.String(), " (will escape cache)")
		}else {
			fmt.Println("Endpoint to test:", _url.String())
		}

		fmt.Println()
	}
}
