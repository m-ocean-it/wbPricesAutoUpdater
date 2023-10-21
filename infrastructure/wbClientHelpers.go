package infrastructure

import (
	"fmt"
	"io"
	"net/http"
)

func nonOkErrorFromResponse(response *http.Response) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("received non-OK status code (%d\n) and could not read the body", response.StatusCode)
	}
	return fmt.Errorf("received non-OK status code (%d\n). body: %s", response.StatusCode, string(body))
}
