package lewd

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/parnurzeal/gorequest"
	"html"
	"strconv"
	"strings"
)

const (
	apiURL = "http://g.e-hentai.org/api.php"
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

	// Make json struct
	jsonStruct := ehentaiRequest{
		Method: method,
	}

	switch method {
	case "gdata":
		{
			jsonStruct.Gidlist = list
			jsonStruct.Namespace = 1
		}
	case "gtoken":
		{
			jsonStruct.Pagelist = list
		}
	}

	// Convert json object to string
	jsonString, err := json.Marshal(jsonStruct)
	if err != nil {
		fmt.Println(err)
	}

	// Post the request
	_, reply, errs := gorequest.New().Post(apiURL).Send(string(jsonString)).EndBytes()

	if err != nil {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	var response ehentaiResponse

	if err := json.Unmarshal(reply, &response); err != nil {
		fmt.Println(err.Error())
	}

	switch method {
	case "gdata":
		{
			return response.Gmetadata
		}
	case "gtoken":
		{
			return response.Tokenlist
		}
	}

	return []galleryMetadata{}
}

func getGalleryTokens(pagelist [][]string) [][]string {
	tokenList := makeRequest("gtoken", pagelist)

	galleries := [][]string{}

	for _, gallery := range tokenList {
		if len(gallery.Error) > 0 {
			fmt.Printf("gid: %d error: %s", gallery.Gid, gallery.Error)
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

func parseGalleryMetadata(s *discordgo.Session, m *discordgo.MessageCreate, galleries []galleryMetadata) {

	if len(m.Embeds) > 0 {
		for _, e := range m.Embeds {
			fmt.Println(e.Type, e.URL)
		}

		// Remove remove all embeds for now
		_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.ID, &discordgo.MessageEmbed{})
		if err != nil {
			fmt.Println(err.Error())
		}

	}

	for _, gallery := range galleries {
		if len(gallery.Error) > 0 {
			fmt.Printf("gid: %d error: %s", gallery.Gid, gallery.Error)
			continue
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

		var fields []*discordgo.MessageEmbedField

		for _, group := range keys {
			var text string
			for i, tag := range tags[group] {
				if i < len(tags[group])-1 {
					text += fmt.Sprintf("%s, ", tag)
				} else {
					text += fmt.Sprintf("%s", tag)
				}
			}

			fields = append(fields, &discordgo.MessageEmbedField{Name: group, Value: text, Inline: true})
		}

		message := discordgo.MessageEmbed{
			URL:       fmt.Sprintf("https://exhentai.org/g/%d/%s/", gallery.Gid, gallery.Token),
			Title:     html.UnescapeString(gallery.Title),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: gallery.Thumb},
			Fields:    fields,
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
