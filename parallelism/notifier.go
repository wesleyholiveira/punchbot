package parallelism

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/cnf/structhash"

	"github.com/wesleyholiveira/punchbot/redis"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	fb "github.com/huandu/facebook"
	log "github.com/sirupsen/logrus"
	"github.com/wesleyholiveira/punchbot/configs"
	"github.com/wesleyholiveira/punchbot/helpers"
	"github.com/wesleyholiveira/punchbot/models"
	"github.com/wesleyholiveira/punchbot/services/punch"
	t "github.com/wesleyholiveira/punchbot/services/twitter"
)

var msgID map[string]string

var channels map[string]string
var mirrors map[string]string

func init() {
	msgID = make(map[string]string)
	mirrors = make(map[string]string)
	channels = helpers.ParseChannels(configs.NotificationChannelsID)

	mirrors["stream"] = configs.PunchEndpoint + "/download-stream/"
	mirrors["zippyshare"] = configs.PunchEndpoint + "/download-zippyshare/"
	mirrors["openload"] = configs.PunchEndpoint + "/download-openload/"
}

func Notifier(s *discordgo.Session, projects chan *[]models.Project) {
	block := false
	for p := range projects {
		var guildID string
		twitter := t.GetClient()
		face := models.GetFacebook()
		prev := models.GetProjects()

		log.Info("Notifier is on")

		for key, tag := range channels {

			ch, _ := s.Channel(key)

			if ch != nil {
				guild, _ := s.Guild(ch.GuildID)
				guildID = guild.ID

				userMention := ""
				if tag != "" {
					username := tag[1:]
					if username != "everyone" {
						for _, m := range guild.Members {
							if m.User.Username == username {
								userMention = m.User.Mention()
								break
							}
						}
						if userMention == "" {
							for _, r := range guild.Roles {
								if r.Name == username && r.Mentionable {
									userMention = fmt.Sprintf("<@&%s>", r.ID)
									break
								}
							}
						}
					} else {
						userMention = "@" + username
					}
					userMention += " "
				}
				go notify(s, twitter, face, p, prev, ch.ID, userMention, block)
			}

			block = true
		}

		myNotifications := models.GetNotifyUser()
		if myNotifications != nil {
			guild, _ := s.Guild(guildID)
			userMention := ""
			for key := range myNotifications {
				ch, _ := s.Channel(key)
				myNots := myNotifications[key]

				if ch != nil {
					if guild != nil {
						if ch.Type == discordgo.ChannelTypeDM {
							if myNots != nil {
								for _, m := range guild.Members {
									if myNots.UserID == m.User.ID {
										log.Info("User found in punch's server")
										for _, userRoleID := range m.Roles {
											for _, role := range guild.Roles {
												if role.Name == "VIP" && role.ID == userRoleID {
													log.Info("The user is a vip!!")
													log.Info("Sending notifications (if exists)")

													myNots.VIP = true
												}
											}
										}

										if _, err := notifyUser(s, p, myNots, key, userMention); err != nil {
											log.Error(err)
										}
										break
									}
								}
							}
						}
					}
				}
			}
		}
		block = false
	}
}

func notify(s *discordgo.Session, t *twitter.Client, f *models.Facebook, current *[]models.Project, prev *[]models.Project, channelID, userMention string, block bool) (bool, error) {
	punchReleases := models.GetProjects()
	cLen := len(*current)
	pLen := len(*prev)
	diff := int(math.Abs(float64(cLen) - float64(pLen)))

	if diff == 0 {
		diff = 1
	}

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d", diff, pLen, cLen)

	currentSlice := (*current)[0:diff]

	for i, c := range currentSlice {
		if !c.AlreadyReleased {
			for _, p := range *prev {
				if c.IDProject != p.IDProject {
					log.Info("PROJECT MATCHED!")

					pCurrent := &c
					sendMessage(s, pCurrent, &p, channelID, userMention)

					if !block {
						sendMessageTwitter(t, pCurrent, channelID)
						sendMessageFacebook(f, pCurrent, channelID)
					}

					(*current)[i].AlreadyReleased = true
					break
				}
			}
		} else {
			log.Warnf("A previosly released project %s[%s] was already released.",
				c.Project, c.HashID)
			log.Warn("Ignoring...")
			break
		}
	}

	*punchReleases = *current

	return true, nil
}

func notifyUser(s *discordgo.Session, current *[]models.Project, myNots *models.Notify, channelID, userMention string) (bool, error) {
	prev := myNots.Projects
	punchReleases := models.GetProjects()
	cLen := len(*current)
	pLen := len(*prev)
	diff := int(math.Abs(float64(cLen) - float64(pLen)))

	if diff == 0 {
		diff = 1
	}

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d", diff, pLen, cLen)
	currentSlice := (*current)[0:diff]

	for i, c := range currentSlice {
		if !c.AlreadyReleased {
			for _, p := range *prev {
				if c.IDProject == p.IDProject {
					log.Info("PROJECT MATCHED! [USER]")

					if !myNots.VIP {
						user, _ := s.User(myNots.UserID)
						s.ChannelMessageSend(channelID,
							fmt.Sprintf("Vejo que você possui interesse em receber notificações "+
								"via DM mas isto é **exclusivo** para usuários **VIP** no Discord.\n"+
								"Leia o canal #regras ou acesse: %s para adquirir seu **VIP!**",
								configs.PunchEndpoint))
						return false, errors.New(fmt.Sprintf("%s Isn't a vip", user.Username))
					}

					sendMessage(s, &c, &p, channelID, userMention)
					(*current)[i].AlreadyReleased = true
					break
				}
			}
		} else {
			log.Warnf("A previosly released project %s[%s] was already released.",
				c.Project, c.HashID)
			log.Warn("Ignoring...")
			break
		}
	}

	*punchReleases = *current

	return true, nil
}

func sendMessageTwitter(t *twitter.Client, c *models.Project, channelID string) {
	msg := fmt.Sprintf("O %s do anime %s acabou de ser lançado! -> %s\n",
		c.Number,
		c.Project,
		configs.PunchEndpoint+c.Link)

	_, _, err := t.Statuses.Update(msg, nil)

	if err != nil {
		log.Error(err)
	}
}

func sendMessageFacebook(f *models.Facebook, c *models.Project, channelID string) {
	re := redis.GetClient()

	data, err := re.Get("fbWhiteList").Result()
	if err != nil {
		log.Error("Facebook notifier: ", err)
		return
	} else {
		projects := make([]models.Project, 0)

		err := json.Unmarshal([]byte(data), &projects)
		if err != nil {
			log.Error("Facebook notifier: ", err)
			return
		}

		for _, project := range projects {
			p := *c
			if project.IDProject == p.IDProject {
				r, _, err := getImage(c)
				fs := f.Session

				if err == nil {
					url := r.Request.URL.String()
					log.Infof("Getting image to facebook: %s", url)

					msg := fmt.Sprintf("O %s do anime %s acabou de ser lançado! -> %s\n",
						c.Number,
						c.Project,
						configs.PunchEndpoint+c.Link)

					log.Info("Publishing the image at the page")
					r, err := fs.Post(fmt.Sprintf("/%s/photos", f.PageID), fb.Params{
						"access_token": fs.AccessToken(),
						"url":          url,
						"is_hidden":    "true",
					})

					if err == nil {
						id := r.GetField("id")
						postID := r.GetField("post_id")

						if err != nil {
							log.Error(err)
							return
						}

						log.Infof("Hidding the image of the feed %s", postID)
						_, err = fs.Post(fmt.Sprintf("/%s/", postID), fb.Params{
							"access_token":        fs.AccessToken(),
							"timeline_visibility": "hidden",
						})

						if err != nil {
							log.Error("Hidding error: ", err)
						}

						log.Info("Publishing the post at the page")
						_, err = fs.Post("/feed", fb.Params{
							"access_token":      fs.AccessToken(),
							"message":           msg,
							"object_attachment": id,
						})

						if err != nil {
							log.Error(err)
						}

						return
					}

					log.Error(err)
					return
				}
				break
			} else {
				log.Warnf("The %s isn't a hyped anime. Ignoring...\n", p.Project)
			}
		}

		log.Error(err)
		return
	}
}

func getImage(c *models.Project) (*http.Response, string, error) {
	screen := c.Screen
	img := strings.Split(screen, "/")
	imgName := img[len(img)-1]

	if !strings.Contains(c.Screen, "http") {
		screen = configs.PunchEndpoint + c.Screen
	}

	httpImage, err := punch.Get(screen)
	return httpImage, imgName, err
}

func sendMessage(s *discordgo.Session, c *models.Project, p *models.Project, channelID, userMention string) (bool, error) {
	r, _, err := getImage(c)
	field := &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("**Novo episódio**"),
		Value: fmt.Sprintf("%s", c.Number),
	}

	link := configs.PunchEndpoint + c.Link
	icon := fmt.Sprintf("%s/imagens/favicon-96x96.png", configs.PunchEndpoint)
	prevHash, _ := structhash.Hash(p.ExtraInfos, 0)
	currentHash, _ := structhash.Hash(c.ExtraInfos, 0)

	arrayFields := make([]*discordgo.MessageEmbedField, 0)
	arrayFields = append(arrayFields, field)

	for _, info := range c.ExtraInfos {
		name := fmt.Sprintf("**%s - %s mb**", strings.ToUpper(info.Format), info.Size)
		values := ""

		values += fmt.Sprintf("[%s](%s)", strings.ToUpper("stream"), mirrors["stream"]+info.ID) + " | "
		values += fmt.Sprintf("[%s](%s)", strings.ToUpper("zippyshare"), mirrors["zippyshare"]+info.ID) + " | "
		values += fmt.Sprintf("[%s](%s)", strings.ToUpper("openload"), mirrors["openload"]+info.ID) + " | "
		values = strings.TrimSuffix(values, " | ")

		extraFields := &discordgo.MessageEmbedField{
			Name:  name,
			Value: values,
		}
		arrayFields = append(arrayFields, extraFields)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "PUNCH! Fansubs",
			IconURL: icon,
		},
		Title:       fmt.Sprintf("%s", c.Project),
		Description: fmt.Sprintf("%s", c.Description),
		URL:         link,
		Color:       65280,
		Fields:      arrayFields,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "PUNCH! Fansubs",
			IconURL: icon,
		},
	}

	if r != nil {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: r.Request.URL.String(),
		}
	}

	if err == nil {
		respImage := r.Body
		defer respImage.Close()

		msg := new(discordgo.Message)

		ch := channelID + c.ID
		if msgID[ch] != "" {
			log.Infof("ExtraInfos: %s,%s", prevHash, currentHash)
			if prevHash != currentHash {
				log.Warn("Editing the message embed")
				msg, err = s.ChannelMessageEditEmbed(channelID, msgID[ch], embed)

				if err != nil {
					log.Error("Error edit embed: ", err)
				} else {
					msgID[ch] = msg.ID
				}
			}
		} else {
			if userMention != "" {
				_, err = s.ChannelMessageSend(channelID, userMention)
				if err != nil {
					log.Error(err)
				}
			}

			msg, err = s.ChannelMessageSendEmbed(channelID, embed)

			if err != nil {
				log.Error("Error create embed: ", err)
			} else {
				msgID[ch] = msg.ID
			}
		}
	} else {
		log.Errorf("[%s] Image error: %s", c.Project, err)
	}

	return true, nil
}
