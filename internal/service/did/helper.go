package did

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReadErrorResponse(res *http.Response) error {
	var data map[string]interface{}

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	if msg, ok := data["message"].(string); ok {
		return fmt.Errorf(msg)
	}

	return fmt.Errorf("unknown error")
}
