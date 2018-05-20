package cli

import (
	"sync"
	"fmt"
)

func StartCommandLineInterface(wg *sync.WaitGroup) {
	fmt.Println("Command Line Interface Starting Up!")
	defer wg.Done()

}
