package helper

import "encoding/json"

func To[T any](v any) (T, error) {
	var t T

	b, err := json.Marshal(v)
	if err != nil {
		return t, err
	}

	if err := json.Unmarshal(b, &t); err != nil {
		return t, err
	}

	return t, nil
}
