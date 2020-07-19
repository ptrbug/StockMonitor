package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func fetch(remoteURL string, queryValues url.Values, cookies []*http.Cookie) ([]byte, error) {

	client := &http.Client{}
	uri, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Host", uri.Host)
	reqest.Header.Add("Referer", uri.String())
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	if cookies != nil {
		for _, v := range cookies {
			reqest.AddCookie(v)
		}
	}

	response, err := client.Do(reqest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		var body []byte
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			gr, err := gzip.NewReader(response.Body)
			if err != nil {
				return nil, err
			}
			body, err = ioutil.ReadAll(gr)
			if err != nil {
				return nil, err
			}

		default:
			body, err = ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}

		}
		return body, nil
	}

	xx, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(xx))

	return nil, fmt.Errorf("response.StatusCode code:%v", response.StatusCode)
}
