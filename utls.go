package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetch(remoteURL string, queryValues url.Values) (body []byte, err error) {

	client := &http.Client{}
	body = nil
	uri, err := url.Parse(remoteURL)
	if err != nil {
		return
	}
	if queryValues != nil {
		values := uri.Query()
		if values != nil {
			for k, v := range values {
				queryValues[k] = v
			}
		}
		uri.RawQuery = queryValues.Encode()
	}
	reqest, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return
	}
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Host", uri.Host)
	reqest.Header.Add("Referer", uri.String())
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	response, err := client.Do(reqest)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(response.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)

				if err != nil && err != io.EOF {
					panic(err)
				}

				if n == 0 {
					break
				}
				body = append(body, buf...)
			}
		default:
			body, _ = ioutil.ReadAll(response.Body)

		}
	}
	return
}
