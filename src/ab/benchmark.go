package ab

import (
	"time"
	"net/http"
	"io/ioutil"
	"fmt"
	"bufio"
	"bytes"
	"sync"
	"strconv"
	"os"
)

type result struct {
	allRequestCount        int
	succeedCount           int
	failedCount            int
	non2xxResponseCount    int
	transferRate           float64
	totalBodySent          int64
	totalTransferred       int64
	htmlTransferred        int64
	requestPerSecond       float64
	timePerRequest         float64
	totalSpend             int64
}

type perResult struct {
	transferred       int64
	bodySent          int64
	htmlTransferred   int64
	timeSpent         int64
	isSucceed         bool
	isNon2xx          bool
	failedReason      error
}

func headerSize(header http.Header) int {
	buffer := bytes.NewBuffer(make([]byte, 300))
	w := bufio.NewWriter(buffer)
	header.WriteSubset(w, nil)
	w.Flush()

	return len(buffer.Bytes())
}

func request() *perResult {
	result := new(perResult)
	result.isSucceed = false

	defer func() {
		if r := recover(); r != nil {
			result.failedReason = fmt.Errorf("%v", r)
		}
	}()

	client := new(http.Client)
	request := new(http.Request)

	_url := *Cfg.EndPointParsed
	request.URL = &_url
	if Cfg.EscapeCache != "" {
		q := request.URL.Query()
		q.Add(Cfg.EscapeCache, strconv.FormatInt(time.Now().UnixNano(), 10))
		request.URL.RawQuery = q.Encode()
	}

	if len(Cfg.Headers) > 0 {
		request.Header = make(http.Header)

		for _, v := range Cfg.Headers {
			request.Header.Set(v[0], v[1])
		}
	}

	if Cfg.Post != "" {
		request.Method = "POST"
		if f, err := os.OpenFile(string(Cfg.Post), os.O_RDONLY, os.FileMode(0666)); err != nil {
			result.failedReason = err

			return result
		}else {
			if fi, err := f.Stat(); err != nil {
				result.failedReason = err
				return result
			}else {
				request.ContentLength = fi.Size()
				request.Body = f
				request.Header.Set("Content-Type", Cfg.ContentType)
			}
		}
	}else {
		request.Method = "GET"
	}

	if len(Cfg.Cookies) > 0 {
		for _, v := range Cfg.Cookies {
			c := new(http.Cookie)
			c.Name = v[0]
			c.Value = v[1]
			request.AddCookie(c)
		}
	}

	timeStart := time.Now().UnixNano()
	response, err := client.Do(request)
	timeEnd := time.Now().UnixNano()

	result.timeSpent = (timeEnd-timeStart)/1e6 // to ms

	if err != nil {
		result.failedReason = err
		return result
	}

	if response.StatusCode != 200 {
		result.isNon2xx = true
	}

	if body, err := ioutil.ReadAll(response.Body); err != nil {
		result.failedReason = err
		return result
	}else {
		result.htmlTransferred = int64(len(body))
		result.transferred = result.htmlTransferred+int64(headerSize(response.Header))
		result.isSucceed = true

		return result
	}
}

func doConcurrence(c int, times int, perResultChan chan <-*perResult) {
	wg := new(sync.WaitGroup)
	for i := 0; i < c; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			perResultChan<-request()
		}()
	}

	wg.Wait()
	fmt.Println("Complete ", Cfg.Concurrence*times, " requests")
}

func Do(done chan int) {
	perResultChan := make(chan *perResult, Cfg.Number)
	totalResult := new(result)

	go func() {
		for {
			r, more := <-perResultChan
			if more {
				if r.isSucceed {
					totalResult.succeedCount++
				}else {
					fmt.Println(r.failedReason)
					totalResult.failedCount++
				}

				if r.isNon2xx {
					totalResult.non2xxResponseCount++
				}

				totalResult.allRequestCount = Cfg.Number
				totalResult.htmlTransferred = totalResult.htmlTransferred+r.htmlTransferred
				totalResult.totalBodySent = totalResult.totalBodySent+r.bodySent
				totalResult.totalTransferred = totalResult.totalTransferred+r.transferred

				totalResult.totalSpend = totalResult.totalSpend+r.timeSpent
				if totalResult.totalSpend > 0 {
					totalResult.requestPerSecond = float64(totalResult.allRequestCount)/float64(totalResult.totalSpend)*1000
				}
			} else {

				fmt.Println()

				fmt.Println("All request:", totalResult.allRequestCount)
				fmt.Println("Time taken:", float64(totalResult.totalSpend)/float64(1000), " [second]")
				fmt.Println("Succeed requests:", totalResult.succeedCount)
				fmt.Println("Failed requests:", totalResult.failedCount)
				fmt.Println("Non2xx requests:", totalResult.non2xxResponseCount)
				fmt.Println("Body sent:", totalResult.totalBodySent, " [bytes]")
				fmt.Println("HTML transferred:", totalResult.htmlTransferred, " [bytes]")
				fmt.Println("Request per second:", fmt.Sprintf("%.2f", totalResult.requestPerSecond), " [#/sec] (mean)")

				fmt.Println()
				done <- 1
				return
			}
		}
	}()

	go func() {
		times := 1
		remain := Cfg.Number;
		c := Cfg.Concurrence

		for {
			if remain == 0 {
				break
			}

			if remain > Cfg.Concurrence {
				c = Cfg.Concurrence
			}else {
				c = remain
			}

			doConcurrence(c, times, perResultChan)
			remain = remain-Cfg.Concurrence
			times++
		}

		close(perResultChan)
	}()
}
