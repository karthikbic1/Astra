package httputils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(ctx context.Context, url string, query map[string]interface{}, headers ...map[string]string) (int, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}

	return perform(ctx, req, headers...)
}

func perform(ctx context.Context, req *http.Request, headers ...map[string]string) (status int, body []byte, err error) {
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return resp.StatusCode, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}

	// We will close only when no error as response is always closed in error case
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusNoContent, nil, err
	}

	return resp.StatusCode, b, nil
}
