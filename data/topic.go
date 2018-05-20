package data

import "time"

type Topic struct {
	ID       int
	Title    string
	Source   string
	SourceID string
	Score    int

	CreatedDate time.Time

	DiscordChannel Channel

	LifecycleState int
}

func AddTopic(top Topic) int {
	return dbInsertTopic(top)
}

func GetActiveTopics() []Topic {
	return dbScrapeTopics(0)
}

func GetAllTopics(limit int) []Topic {
	return dbScrapeTopics(limit)
}
