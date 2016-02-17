package lewd

import (
	"github.com/softashell/lewdbot-discord/regex"
)

func ParseLinks(text string) (bool, string) {
	ex_galleries := [][]string{} // id, token
	ex_pages := [][]string{}     // id, page_token, page_number

	// exhentai
	ex_gallery_links := regex.ExGalleryLink.FindAllStringSubmatch(text, -1)
	ex_gallery_page_links := regex.ExGalleryPage.FindAllStringSubmatch(text, -1)

	nh_galleries := []string{}

	// nhentai
	nh_gallery_links := regex.NhGalleryLink.FindAllStringSubmatch(text, -1)

	for _, link := range ex_gallery_links {
		id := link[1]
		token := link[2]

		ex_galleries = append(ex_galleries, []string{id, token})
	}

	for _, link := range ex_gallery_page_links {
		page_token := link[1]
		id := link[2]
		page_number := link[3]

		ex_pages = append(ex_pages, []string{id, page_token, page_number})
	}

	if len(ex_pages) > 0 {
		for _, gallery := range get_gallery_tokens(ex_pages) {
			ex_galleries = append(ex_galleries, gallery)
		}
	}

	// Doesn't actually do anything with it yet, maybe later
	for _, link := range nh_gallery_links {
		id := link[1]

		nh_galleries = append(nh_galleries, id)
	}

	var reply string

	if len(ex_galleries) > 0 {
		gallery_metadata := get_gallery_metadata(ex_galleries)
		reply = parse_gallery_metadata(gallery_metadata)
	} else if len(nh_galleries) > 0 {
		reply = "```css\n>nhentai\n```"
	} else {
		// Didn't find anything
		return false, reply
	}

	return true, reply
}
