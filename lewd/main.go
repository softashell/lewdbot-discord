package lewd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/softashell/lewdbot-discord/regex"
)

// ParseLinks Returns gallery metadata from founds links in input
func ParseLinks(s *discordgo.Session, m *discordgo.MessageCreate, text string) bool {
	exGalleries := [][]string{} // id, token
	exPages := [][]string{}     // id, page_token, page_number

	// exhentai
	exGalleryLinks := regex.ExGalleryLink.FindAllStringSubmatch(text, -1)
	exGalleryPageLinks := regex.ExGalleryPage.FindAllStringSubmatch(text, -1)

	nhGalleries := []string{}

	// nhentai
	nhGalleryLinks := regex.NhGalleryLink.FindAllStringSubmatch(text, -1)

	fmt.Println(len(m.Embeds))

	for _, e := range m.Embeds {
		fmt.Println(e)
	}

	for _, link := range exGalleryLinks {
		id := link[1]
		token := link[2]

		exGalleries = append(exGalleries, []string{id, token})
	}

	for _, link := range exGalleryPageLinks {
		pageToken := link[1]
		id := link[2]
		pageNumber := link[3]

		exPages = append(exPages, []string{id, pageToken, pageNumber})
	}

	if len(exPages) > 0 {
		for _, gallery := range getGalleryTokens(exPages) {
			exGalleries = append(exGalleries, gallery)
		}
	}

	// Doesn't actually do anything with it yet, maybe later
	for _, link := range nhGalleryLinks {
		id := link[1]

		nhGalleries = append(nhGalleries, id)
	}

	if len(exGalleries) > 0 {
		galleryMetadata := getGalleryMetadata(exGalleries)
		parseGalleryMetadata(s, m, galleryMetadata)
	} else if len(nhGalleries) > 0 {
		reply := "```css\n>nhentai\n```"
		_, err := s.ChannelMessageSend(m.ChannelID, reply)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		// Didn't find anything
		return false
	}

	return true
}
