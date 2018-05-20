package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jdanek4/Unite-Backend/config"
	"github.com/jdanek4/Unite-Backend/data"
	"fmt"
	"os"
)

func SetupBot() {
	var err error
	discordSession, err = discordgo.New("Bot " + config.Config().Discord_token)
	if err != nil {
		fmt.Println("Error Connecting to Bot!")
	}
}

func resetDiscordServer() {
	discord := getDiscordBot()

	err := discord.Open()

	if err != nil {
		fmt.Println("Error: Connecting to discord Bot")
		os.Exit(1)
	}

	user, err := discord.User("@me")
	if err != nil {
		fmt.Println("Error Connecting to Bot!")
		return
	}
	discord.State.User = user

	guild := discord.State.Guilds[0]

	// Purge All Channels that are currently active and set their state to Matured
	for _, channel := range guild.Channels {
		if channel.ParentID != "446408325786238976" && channel.ID != "446408325786238976" {
			discord.ChannelDelete(channel.ID)
		}
	}

	for len(guild.Channels) > 2 {
		for _, channel := range guild.Channels {
			if channel.ParentID != "446408325786238976" && channel.ID != "446408325786238976" {
				discord.ChannelDelete(channel.ID)
			}
		}
	}

	discord.Close()

}

func getDiscordBot() (*discordgo.Session) {
	return discordSession
}

func removeChannelHierarchy(chanID string, discord *discordgo.Session) {
	dbChannels := data.GetChannels()
	for _, chans := range dbChannels {
		if chans.ChannelID == chanID {

			discord.ChannelDelete(chans.ChatID)
			discord.ChannelDelete(chans.VoiceID)
			discord.ChannelDelete(chans.ChannelID)
			data.RemoveChannel(chans.ID)

		}
	}
}

const (
	ChannelTypeGuildText     = "0"
	ChannelTypeDM            = "1"
	ChannelTypeGuildVoice    = "2"
	ChannelTypeGroupDM       = "3"
	ChannelTypeGuildCategory = "4"
)

func addChannelHierarchyForTopic(topic data.Topic, discord *discordgo.Session, guild string) {
	topCat, err := discord.GuildChannelCreate(guild, topic.Title, ChannelTypeGuildCategory)
	topVoice, err := discord.GuildChannelCreate(guild, "voice", ChannelTypeGuildCategory)
	topText, err := discord.GuildChannelCreate(guild, "chat", ChannelTypeGuildText)

	topVoice, err = discord.ChannelEditComplex(topVoice.ID, &discordgo.ChannelEdit{ParentID: topCat.ID})
	topText, err = discord.ChannelEditComplex(topText.ID, &discordgo.ChannelEdit{ParentID: topCat.ID})

	if err != nil {
		return
	}

	discord.State.ChannelAdd(topCat)
	discord.State.ChannelAdd(topVoice)
	discord.State.ChannelAdd(topText)

	var channel data.Channel
	channel.ChannelID = topCat.ID
	channel.VoiceID = topVoice.ID
	channel.ChatID = topText.ID
	channel.Topic = topic.ID
	data.AddChannel(channel)

	discord.ChannelMessageSend(topText.ID, topic.Title)
	discord.ChannelMessageSend(topText.ID, "----------------------------------------------------------")
	discord.ChannelMessageSend(topText.ID, "This lab will expires in 1 hour or after prolonged inactivity!")

	inviteParams := discordgo.Invite{MaxAge: 14400, MaxUses: 0, Temporary: false}

	chatInvite, err := discord.ChannelInviteCreate(channel.ChatID, inviteParams)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	voiceInvite, err := discord.ChannelInviteCreate(channel.VoiceID, inviteParams)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	channel.VoiceInvite = voiceInvite.Code
	channel.VoiceInvite = chatInvite.Code

	data.AddChannel(channel)
}
