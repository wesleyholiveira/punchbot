package parallelism

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

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
var Channels map[string]string
var mirrors map[string]string

func init() {
	msgID = make(map[string]string)
	mirrors = make(map[string]string)
	Channels = helpers.ParseChannels(configs.NotificationChannelsID)

	mirrors["stream"] = configs.PunchEndpoint + "/download-stream/"
	mirrors["zippyshare"] = configs.PunchEndpoint + "/download-zippyshare/"
	mirrors["openload"] = configs.PunchEndpoint + "/download-openload/"
}

func Notifier(s *discordgo.Session, projects chan *[]models.Project) {
	block := false
	for p := range projects {
		twitter := t.GetClient()
		face := models.GetFacebook()
		prev := models.GetProjects()

		log.Infof("Notifier is on")

		for key, tag := range Channels {

			ch, _ := s.Channel(key)

			if ch != nil {
				guild, _ := s.Guild(ch.GuildID)

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

				log.Infof("Notifing for #%s channel", ch.Name)
				notify(s, twitter, face, p, prev, ch.ID, userMention, block)
			}

			block = true
		}

		myNotifications := models.GetNotifyUser()
		if myNotifications != nil {
			userMention := ""
			for key, myNots := range myNotifications {
				if _, err := notifyUser(s, p, myNots, key, userMention); err != nil {
					log.Error(err)
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

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d (GLOBAL)", diff, pLen, cLen)

	currentSlice := (*current)[0:diff]
	prevSlice := (*current)[1:diff]

	for i, c := range currentSlice {
		if !c.AlreadyReleased {
			for _, p := range prevSlice {
				if c.IDProject != p.IDProject {
					log.Info("PROJECT MATCHED!")

					sendMessage(s, c, p, channelID, userMention)

					if !block {
						sendMessageTwitter(t, &c, channelID)
						sendMessageFacebook(f, &c, channelID)
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
	red := redis.GetClient()
	prev := myNots.Projects
	punchReleases := models.GetProjects()
	cLen := len(*current)
	pLen := len(*prev)
	diff := pLen

	if diff == 0 {
		diff = 1
	}

	log.Infof("Diff: %d, PREV PROJECTS: %d, CURRENT PROJECTS: %d (USER)", diff, pLen, cLen)
	currentSlice := (*current)[0]

	for i, p := range *prev {
		if currentSlice.IDProject == p.IDProject && !p.AlreadyReleased {
			user, _ := s.User(myNots.UserID)
			log.Infof("PROJECT MATCHED! (%s) -> [%s]", user.Username, p.Project)

			if !myNots.VIP {
				return false, fmt.Errorf("%s Isn't a vip", user.Username)
			}

			sendMessage(s, currentSlice, p, channelID, userMention)
			(*current)[0].AlreadyReleased = true
			(*prev)[i].ExtraInfos = currentSlice.ExtraInfos

			if len((*prev)[i].ExtraInfos) == 4 {
				(*prev)[i].AlreadyReleased = true
			}

			notify := &models.Notify{UserID: myNots.UserID, Projects: myNots.Projects}
			notify.VIP = myNots.VIP
			notifyMarshal, err := json.Marshal(notify)

			if err != nil {
				return false, fmt.Errorf("NotifierUser Redis Marshal %s", err)
			}

			notifyMarshalStr := string(notifyMarshal)
			err = red.Set(channelID, notifyMarshalStr, 0).Err()

			if err != nil {
				return false, fmt.Errorf("NotifierUser Redis Set %s", err)
			}

			red.Save().Err()
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
				r, _, err := getImage(p)
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

func getImage(c models.Project) (*http.Response, string, error) {
	screen := c.Screen
	img := strings.Split(screen, "/")
	imgName := img[len(img)-1]

	if !strings.Contains(c.Screen, "http") {
		screen = configs.PunchEndpoint + c.Screen
	}

	httpImage, err := punch.Get(screen)
	return httpImage, imgName, err
}

func sendMessage(s *discordgo.Session, c models.Project, p models.Project, channelID, userMention string) (bool, error) {
	r, _, err := getImage(c)
	field := &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("**Novo episódio**"),
		Value: fmt.Sprintf("%s", c.Number),
	}

	link := configs.PunchEndpoint + c.Link
	icon := fmt.Sprintf("%s/imagens/favicon-96x96.png", configs.PunchEndpoint)
	lenPrev := len(p.ExtraInfos)
	lenCurrent := len(c.ExtraInfos)

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
			log.Infof("ExtraInfos: current: %d, prev: %d", lenCurrent, lenPrev)
			if lenCurrent > 0 && lenCurrent > lenPrev {
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
