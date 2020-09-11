package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func GetAuthorNickname(s *discordgo.Session, m *discordgo.Message) string {
	return GetNickName(s, m.Author.ID, m.GuildID)
}

func GetNickName(s *discordgo.Session, userId string, guildId string) string {
	member, err := s.GuildMember(guildId, userId)
	he(err)
	return member.Nick
}
func amImentionned(s *discordgo.Session, m *discordgo.Message) bool {
	for _, username := range m.Mentions {
		if username.Bot && username.ID == s.State.User.ID {
			return true
		}
	}
	return false
}

func he(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

func hde(_ interface{}, err error) {
	he(err)
}
