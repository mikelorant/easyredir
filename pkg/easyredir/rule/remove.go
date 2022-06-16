package rule

import (
	"fmt"
	"net/http"
)

func RemoveRule(cl ClientAPI, id string) (res bool, err error) {
	pathQuery := buildRemoveRule(id)
	reader, err := cl.SendRequest(pathQuery, http.MethodDelete, nil)
	if err != nil {
		return false, fmt.Errorf("unable to send request: %w", err)
	}
	reader.Close()

	return true, nil
}

func buildRemoveRule(id string) string {
	return fmt.Sprintf("/rules/%v", id)
}
