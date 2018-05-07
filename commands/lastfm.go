package commands

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/parnurzeal/gorequest"
	"github.com/softashell/lewdbot-discord/config"

	"fmt"

	"strings"

	log "github.com/Sirupsen/logrus"
)

const apiURL = "http://ws.audioscrobbler.com/2.0/?"

type lastfmreply struct {
	Error        int    `json:"error"`
	Message      string `json:"message"`
	Recenttracks struct {
		Track []struct {
			Name   string `json:"name"`
			Artist struct {
				Text string `json:"#text"`
			} `json:"artist"`
			Album struct {
				Text string `json:"#text"`
			} `json:"album"`
			Attr struct {
				Nowplaying string `json:"nowplaying"`
			} `json:"@attr,omitempty"`
			Date struct {
				Uts string `json:"uts"`
			} `json:"date,omitempty"`
		} `json:"track"`
	} `json:"recenttracks"`
}

func registerLastfmProfile(UserID string, arg string) string {
	arg = strings.TrimSpace(arg)

	if len(arg) > 15 || len(arg) < 2 {
		return "Are you trying to trick me?"
	}

	config.SetLastfmUsername(UserID, arg)

	return fmt.Sprintf("Changed your last.fm username to %q", arg)
}

func removeLastfmProfile(UserID string) string {
	username, err := config.GetLastfmUsername(UserID)
	if err != nil {
		log.Errorf("removeLastfmProfile >> %v", err)
		return "Who the HELL are you?~"
	}

	config.RemoveLastfmUsername(UserID)

	return fmt.Sprintf("Removed your saved last.fm username (%s)~", username)
}

func spamNowPlayingUser(UserID string) string {
	username, err := config.GetLastfmUsername(UserID)
	if err != nil {
		log.Errorf("spamNowPlayingUser >> %v", err)
		return "You haven't registered your last.fm profile yet! Use ``!np set username`` to register~"
	}

	np, err := getNowPlaying(UserID, username)
	if err != nil {
		log.Errorf("spamNowPlayingUser >> %v", err)

		np = "Maybe you should try playing something~"
	}

	return "```" + np + "```"
}

func spamNowPlayingServer(s *discordgo.Session, GuildID string) string {
	g, err := s.State.Guild(GuildID)
	if err != nil {
		log.Errorf("s.State.Guild >> %v", err)
		return "You fucking broke it~"
	}

	var out string

	for _, m := range g.Members {
		username, err := config.GetLastfmUsername(m.User.ID)
		if err != nil {
			continue
		}

		np, err := getNowPlaying(m.User.ID, username)
		if err != nil {
			log.Warn(err)
			continue
		}

		out += m.User.Username + ": " + np + "\n"
	}

	if len(out) < 1 {
		return "Nothing to see here~"
	}

	return "```" + out + "```"
}

func getNowPlaying(UserID string, username string) (string, error) {
	params := fmt.Sprintf("method=user.getRecentTracks&user=%s&api_key=%s&limit=1&format=json", username, config.GetLastfmKey())
	url := apiURL + params

	var response lastfmreply

	// Post the request
	_, reply, errs := gorequest.New().Get(url).Timeout(5 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)

		return "", fmt.Errorf("You fucking broke it")
	}

	if response.Error > 0 {
		log.Error(response.Error, response.Message)

		// 6 : Invalid parameters - Your request is missing a required parameter
		if response.Error == 6 && response.Message == "User not found" {
			log.Errorf("getNowPlaying >> Removing username %q from user %q since last.fm tells us it doesn't exist.", username, UserID)
			config.RemoveLastfmUsername(UserID)
		}

		return "", fmt.Errorf("%d %s", response.Error, response.Message)
	}

	if len(response.Recenttracks.Track) < 1 {
		return "", fmt.Errorf("Didn't get any tracks from this nerd")
	}

	artist := response.Recenttracks.Track[0].Artist.Text
	track := response.Recenttracks.Track[0].Name

	if len(artist) < 1 || len(track) < 1 {
		return "", fmt.Errorf("Empty metadata")
	}

	out := fmt.Sprintf("%s - %s", artist, track)

	if len(response.Recenttracks.Track[0].Attr.Nowplaying) < 1 {
		i, err := strconv.ParseInt(response.Recenttracks.Track[0].Date.Uts, 10, 64)
		if err != nil {
			log.Warning(err)
		} else {
			duration := time.Since(time.Unix(i, 0))

			if duration.Hours() < 6 {
				if duration.Minutes() >= 60 {
					out += fmt.Sprintf(" [%.fh ago]", duration.Hours())
				} else if duration.Seconds() >= 60 {
					out += fmt.Sprintf(" [%.fm ago]", duration.Minutes())
				}
			} else {
				return out, fmt.Errorf("Last track played too long ago")
			}
		}
	}

	return out, nil
}
