package lewd

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/softashell/lewdbot-discord/regex"
	"html"
	"strconv"
	"strings"
)

const (
	api_url = "http://g.e-hentai.org/api.php"
)

type ehentai_request struct {
	Method    string     `json:"method"`
	Gidlist   [][]string `json:"gidlist,omitempty"`
	Pagelist  [][]string `json:"pagelist,omitempty"`
	Namespace int        `json:"namespace,omitempty"`
}

type ehentai_response struct {
	Gmetadata []gallery_metadata `json:"gmetadata"`
	Tokenlist []gallery_metadata `json:"tokenlist"`
}

type gallery_metadata struct {
	Gid   int      `json:"gid"`
	Token string   `json:"token"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	Error string   `json:"error"`
}

func ParseLinks(text string) (bool, string) {
	galleries := [][]string{} // id, token
	pages := [][]string{}     // id, page_token, page_number

	gallery_links := regex.GalleryLink.FindAllStringSubmatch(text, -1)
	gallery_page_links := regex.GalleryPage.FindAllStringSubmatch(text, -1)

	for _, link := range gallery_links {
		id := link[1]
		token := link[2]

		galleries = append(galleries, []string{id, token})
	}

	for _, link := range gallery_page_links {
		page_token := link[1]
		id := link[2]
		page_number := link[3]

		pages = append(pages, []string{id, page_token, page_number})
	}

	if len(pages) > 0 {
		for _, gallery := range get_gallery_tokens(pages) {
			galleries = append(galleries, gallery)
		}
	}

	if len(galleries) < 1 {
		return false, ""
	}

	gallery_metadata := get_gallery_metadata(galleries)

	reply := parse_gallery_metadata(gallery_metadata)

	return true, reply
}

func make_json_request(method string, list [][]string) []gallery_metadata {

	// Make json struct
	jsonStruct := ehentai_request{
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
	_, reply, errs := gorequest.New().Post(api_url).Send(string(jsonString)).EndBytes()

	if err != nil {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}

	var response ehentai_response

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

	return []gallery_metadata{}
}

func get_gallery_tokens(pagelist [][]string) [][]string {
	tokenList := make_json_request("gtoken", pagelist)

	galleries := [][]string{}

	for _, gallery := range tokenList {
		if len(gallery.Error) > 0 {
			fmt.Printf("gid: %s error: %s", gallery.Gid, gallery.Error)
			continue
		}

		galleries = append(galleries, []string{strconv.Itoa(gallery.Gid), gallery.Token})
	}

	return galleries
}

func get_gallery_metadata(galleries [][]string) []gallery_metadata {
	galleryMetadata := make_json_request("gdata", galleries)

	return galleryMetadata
}

func parse_gallery_metadata(galleries []gallery_metadata) string {
	var text string
	var add_url bool

	if len(galleries) > 1 {
		add_url = true
	}

	for _, gallery := range galleries {
		if len(gallery.Error) > 0 {
			fmt.Printf("gid: %s error: %s", gallery.Gid, gallery.Error)
			continue
		}

		text += fmt.Sprintf("**%s**", html.UnescapeString(gallery.Title))

		if add_url {
			// DISCORD A SHIT
			text += fmt.Sprintf(" *exhentai.org/g/%d/%s/*", gallery.Gid, gallery.Token)
		}

		text += fmt.Sprintf("\n```")

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

		for _, group := range keys {
			text += fmt.Sprintf("%s: ", group)
			for i, tag := range tags[group] {
				if i < len(tags[group])-1 {
					text += fmt.Sprintf("%s, ", tag)
				} else {
					text += fmt.Sprintf("%s\n", tag)
				}
			}
		}

		text += fmt.Sprintf("\n```")
	}

	return text
}
