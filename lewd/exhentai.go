package lewd

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

const (
	apiURL = "https://e-hentai.org/api.php"
)

type ehentaiRequest struct {
	Method    string     `json:"method"`
	Gidlist   [][]string `json:"gidlist,omitempty"`
	Pagelist  [][]string `json:"pagelist,omitempty"`
	Namespace int        `json:"namespace,omitempty"`
}

type ehentaiResponse struct {
	Gmetadata []galleryMetadata `json:"gmetadata"`
	Tokenlist []galleryMetadata `json:"tokenlist"`
}

type galleryMetadata struct {
	Gid   int      `json:"gid"`
	Token string   `json:"token"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Thumb string   `json:"thumb"`
	Error string   `json:"error"`
}

func makeRequest(method string, list [][]string) []galleryMetadata {
	jsonStruct := ehentaiRequest{
		Method: method,
	}

	switch method {
	case "gdata":
		jsonStruct.Gidlist = list
		jsonStruct.Namespace = 1
	case "gtoken":
		jsonStruct.Pagelist = list
	}

	var response ehentaiResponse

	// Post the request
	_, reply, errs := gorequest.New().Post(apiURL).SendStruct(jsonStruct).Timeout(10 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)

		return []galleryMetadata{}
	}

	switch method {
	case "gdata":
		return response.Gmetadata
	case "gtoken":
		return response.Tokenlist
	}

	return []galleryMetadata{}
}

func getGalleryTokens(pagelist [][]string) [][]string {
	tokenList := makeRequest("gtoken", pagelist)

	galleries := [][]string{}

	for _, gallery := range tokenList {
		if len(gallery.Error) > 0 {
			log.Warnf("gid: %d error: %s", gallery.Gid, gallery.Error)
			continue
		}

		galleries = append(galleries, []string{strconv.Itoa(gallery.Gid), gallery.Token})
	}

	return galleries
}

func getGalleryMetadata(galleries [][]string) []galleryMetadata {
	galleryMetadata := makeRequest("gdata", galleries)

	return galleryMetadata
}

func parseGalleryMetadata(s *discordgo.Session, channel string, galleries []galleryMetadata) {
	for _, gallery := range galleries {
		if len(gallery.Error) > 0 {
			log.Warnf("gid: %d error: %s", gallery.Gid, gallery.Error)
		}

		var keys []string // Need to keep slice with keys since map doesn't preserve order
		tags := map[string][]string{}

		for _, _tag := range gallery.Tags {
			_tag := strings.Split(_tag, ":")

			group, tag := "misc", ""

			if len(_tag) > 1 { // group:tag_name
				group = _tag[0]
				tag = _tag[1]
			} else { // tag_name
				tag = _tag[0]
			}

			tags[group] = append(tags[group], tag)

			// Only add new key is last one was different
			if len(keys) > 0 && keys[len(keys)-1] == group {
				continue
			}

			keys = append(keys, group)
		}

		var description string
		var fields []*discordgo.MessageEmbedField

		for _, group := range keys {
			var text string
			for i, tag := range tags[group] {
				if i < len(tags[group])-1 {
					text += fmt.Sprintf("%s, ", tag)
				} else {
					text += fmt.Sprintf("%s", tag)
				}

				if((group == "male" && tag == "yaoi") || (group == "female" && tag == "futanari")) {
					description = ":exclamation: :warning: :exclamation: GLOBOHOMO SHIT :exclamation: :warning: :exclamation:"
				}
			}

			fields = append(fields, &discordgo.MessageEmbedField{Name: group, Value: text, Inline: true})
		}

		message := discordgo.MessageEmbed{
			URL:       fmt.Sprintf("https://exhentai.org/g/%d/%s/", gallery.Gid, gallery.Token),
			Title:     html.UnescapeString(gallery.Title),
			Description: description,
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: gallery.Thumb},
			Fields:    fields,
		}

		_, err := s.ChannelMessageSendEmbed(channel, &message)
		if err != nil {
			log.Warn("s.ChannelMessageSendEmbed >>", err)
		}
	}
}
