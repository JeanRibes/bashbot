package main

import "github.com/bwmarrin/discordgo"

func GetAuthorNickname(s *discordgo.Session, m *discordgo.Message) string {
	return GetNickName(s, m.Author.ID, m.GuildID)
}

func GetNickName(s *discordgo.Session, userId string, guildId string) string {
	member, err := s.GuildMember(guildId, userId)
	he(err)
	return member.Nick
}
