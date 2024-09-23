package module

import "bytes"

func parseModuleNames(input []byte) [][]byte {
	if len(input) == 0 {
		return nil
	}

	moduleName := [][]byte{}
	lastNul := 0
	for i := range input {
		if input[i] == 0 {
			word := input[lastNul:i]
			if slshIdx := bytes.IndexByte(word, '/'); slshIdx > -1 {
				word = word[:slshIdx]
			}

			moduleName = append(moduleName, word)
			lastNul = i + 1
		}
	}

	if lastNul != len(input) {
		moduleName = append(moduleName, input[lastNul:])
	}

	return moduleName
}
