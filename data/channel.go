package data

import "time"

type Channel struct {
	ID    int
	Topic int

	ChannelID string

	VoiceID     string
	VoiceInvite string

	ChatID     string
	ChatInvite string

	CreatedDate time.Time
}

func GetChannels() []Channel {
	return dbScrapeChannels()
}
func AddChannel(channel Channel) {
	dbInsertChannel(channel)
}
func RemoveChannel(id int) {
	dbDeleteChannel(id)
}
