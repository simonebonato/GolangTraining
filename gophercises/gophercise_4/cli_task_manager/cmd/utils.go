package cmd

import (
	"fmt"
	"strconv"
)

func argsToIds(args []string) []int {
	var ids []int
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Failed to parse: ", arg)
			return nil
		} else {
			ids = append(ids, id)
		}
	}
	return ids
}
