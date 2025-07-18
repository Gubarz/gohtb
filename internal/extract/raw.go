package extract

import (
	"bytes"
	"io"
	"net/http"
)

func Raw(resp *http.Response) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}
