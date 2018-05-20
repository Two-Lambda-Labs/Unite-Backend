package discord

import (
	"fmt"
	"sync"
	"github.com/bwmarrin/discordgo"
)

var discordSession *discordgo.Session

func StartDiscordManager(wg *sync.WaitGroup) {
	fmt.Println("Discord Overlord Starting Up!")
	defer wg.Done()

	SetupBot()

	resetDiscordServer()

}
