package models

import (
	"github.com/bwmarrin/discordgo"
)

type FnCmd func(projects []Project, s *discordgo.Session, m *discordgo.MessageCreate, args []string)
type Commands map[string]FnCmd
type TNotify map[string]*Notify
