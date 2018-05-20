package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
	"github.com/jdanek4/Antistatic/app/config"
	"strconv"
)

var DataBase *sql.DB

func GetDBConnection() *sql.DB {
	var err error
	if DataBase == nil {
		DataBase, err = sql.Open("mysql", BuildDatasource())
	} else if err = DataBase.Ping(); err != nil {
		DataBase, err = sql.Open("mysql", BuildDatasource())
	}

	if err != nil {
		fmt.Println("Error connecting to Database!")
		return &sql.DB{}
	}
	return DataBase
}

func dbScrapeChannels() []Channel {
	db := GetDBConnection()

	rows, err := db.Query("SELECT id, channel_id, voice_id, chat_id, topic FROM Channels");
	if err != nil {
		return []Channel{}
	}

	var entries []Channel
	defer rows.Close()
	defer db.Close()

	for rows.Next() {
		var ch Channel

		rows.Scan(&ch.ID, &ch.ChannelID, &ch.VoiceID, &ch.ChatID, &ch.Topic)

		entries = append(entries, ch)
	}

	return entries
}
func dbScrapeTopics(limit int) []Topic {
	db := GetDBConnection()

	var err error
	var rows *sql.Rows

	if limit != 0 {
		rows, err = db.Query("SELECT id, title, source, source_id, score, created_date, lifecycle_state FROM Topics ORDER BY score desc LIMIT " + strconv.Itoa(limit) + ";");
		if err != nil {
			return []Topic{}
		}
	} else {
		rows, err = db.Query("SELECT id, title, source, source_id, score, created_date, lifecycle_state FROM Topics ORDER BY score desc;");
		if err != nil {
			return []Topic{}
		}
	}
	var entries []Topic
	defer rows.Close()
	defer db.Close()

	for rows.Next() {
		top := Topic{}

		rows.Scan(&top.ID, &top.Title, &top.Source, &top.SourceID, &top.Score, &top.CreatedDate, &top.LifecycleState)

		entries = append(entries, top)
	}

	return entries
}

func dbInsertChannel(channel Channel) {
	db := GetDBConnection()

	num := db.QueryRow("select count(1) from Channels Where topic=?;", channel.Topic)

	channelExists := 10
	err := num.Scan(&channelExists)
	if err != nil || channelExists != 0 {
		return
	}

	var datetime = time.Now()
	datetime.Format(time.RFC3339)
	stmt, err := db.Prepare("INSERT INTO Channels (topic, channel_id, voice_id, voice_invite, chat_id, chat_invite, created_date) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	_, err = stmt.Exec(channel.Topic, channel.ChannelID, channel.VoiceID, channel.VoiceInvite, channel.ChatID, channel.ChatInvite, datetime)
	if err != nil {
		return
	}

	db.Close()
	return
}

func dbDeleteChannel(id int) {
	db := GetDBConnection()

	stmt, err := db.Prepare("DELETE FROM Channels WHERE id=?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return
	}

	db.Close()
	return
}

//(&top.ID, &top.Title, &top.Content, &top.Score, &top.User, &top.Created_date, &top.Voice_server, &top.Lifecycle_state, &top.Archive_audio, &top.Archive_chat)
func dbInsertTopic(top Topic) int {
	db := GetDBConnection()

	stmt, err := db.Prepare("INSERT INTO Topics (title, source, source_id, score, created_date, lifecycle_state) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return -1
	}
	var datetime = time.Now()
	datetime.Format(time.RFC3339)
	_, err = stmt.Exec(top.Title, top.Source, top.SourceID, top.Score, datetime, top.LifecycleState)
	if err != nil {
		return -1
	}

	stmt2, err := db.Prepare("SELECT LAST_INSERT_ID();")
	if err != nil {
		return -1
	}
	var index int
	err = stmt2.QueryRow().Scan(&index)
	if err != nil {
		return -1
	}
	db.Close()
	return index
}

func BuildDatasource() string {
	datasource := ""

	datasource += config.Config().DB_user
	datasource += ":"
	datasource += config.Config().DB_pass
	datasource += "@tcp("
	datasource += config.Config().DB_ip
	datasource += ")/"
	datasource += config.Config().DB_name
	datasource += "?charset=utf8&parseTime=True"

	return datasource
}
