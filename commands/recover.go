package commands

import "fmt"

type recoverStrategy func()

func RecoverStrategy(command string) recoverStrategy {
	return func() {
		if r := recover(); r != nil {
			fmt.Printf("The request was invalid or cannot be served. For more details, please run: `apihub %s -h`.", command)
		}
	}
}
