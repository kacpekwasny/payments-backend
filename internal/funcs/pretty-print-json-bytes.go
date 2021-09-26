package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PPJsonBytes(b []byte) {
	// pretty-print the json
	go func() {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, b, "", "    ")
		if err != nil {
			fmt.Println("PPJsonBytes(...) err:", err)
			return
		}
		fmt.Println(prettyJSON.String())
	}()
}
