package Unite_Backend

import (
	"fmt"
	"sync"
	"github.com/jdanek4/Unite-Backend/discord"
	"github.com/jdanek4/Unite-Backend/cli"
)

func main() {

	fmt.Println("Starting Unite Backend!")

	var wg sync.WaitGroup

	wg.Add(2)

	go discord.StartDiscordManager(&wg)

	go cli.StartCommandLineInterface(&wg)

	wg.Wait()

	fmt.Println("Shutting down...")
}
