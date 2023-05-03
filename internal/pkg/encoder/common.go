package encoder

import "encoding/json"

func EncodeToJSON(a any) string {
	bytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
