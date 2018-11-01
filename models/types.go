package models

import (
	"github.com/bwmarrin/discordgo"
)

type FnCmd func(s *discordgo.Session, channel string, args []string)
type Commands map[string]FnCmd
type TNotify map[string]*Notify

type GetProjectsType uint32

const (
	Calendar GetProjectsType = iota
	Home
)
