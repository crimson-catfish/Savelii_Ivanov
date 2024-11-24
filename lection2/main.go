package main

import (
	"encoding/json"
	"fmt"
	"os"

	"entrance/lection2/operproc"
)

const (
	envVar = "FILE_TO_PROCESS"
	output = "lection2/out.json"
)

func main() {
	file, err := operproc.GetFilePath(envVar)

	dat, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	var unformattedOperations []map[string]any
	if err := json.Unmarshal(dat, &unformattedOperations); err != nil {
		fmt.Println(err)
	}

	processedOperations, err := operproc.ProcessOperations(unformattedOperations)
	if err != nil {
		return
	}

	processedOperationsJSON, err := json.Marshal(processedOperations)
	if err != nil {
		fmt.Println(err)
	}
	if err := os.WriteFile(output, processedOperationsJSON, 0600); err != nil {
		fmt.Println(err)
	}
}
