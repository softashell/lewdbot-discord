package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func create_client(id string, hash string) *http.Client {
	cookieJar, _ := cookiejar.New(nil)

	cookies := []*http.Cookie{
		{
			Name:   "ipb_member_id",
			Value:  id,
			Path:   "/",
			Domain: ".exhentai.org",
		},
		{
			Name:   "ipb_pass_hash",
			Value:  hash,
			Path:   "/",
			Domain: ".exhentai.org",
		},
	}

	cookieURL, _ := url.Parse("http://exhentai.org")
	cookieJar.SetCookies(cookieURL, cookies)

	client := &http.Client{
		Jar: cookieJar,
	}

	return client
}

func parse_links(text string) (bool, string) {
	links := regexp.MustCompile(`http://(ex|g\.e-)hentai.org/g/[[:alnum:]]+/[[:alnum:]]+`).FindAllString(text, -1)

	found := false
	reply := ""

	for _, link := range links {
		reply += parse_link(link)

		found = true
	}

	return found, reply
}

func parse_link(url string) string {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	title_en := doc.Find("#gn").Text()
	//title_jp := doc.Find("#gj").Text()

	text := fmt.Sprintf("%s\n", title_en)

	taglist := doc.Find("#taglist > table > tbody > tr")

	if taglist.Length() > 0 {
		groups := taglist.Find(".tc")
		for i := range groups.Nodes {
			group := groups.Eq(i)

			text += fmt.Sprintf("%s ", group.Text())

			tags := group.Next().Find(".gt, .gtl")
			for i := range tags.Nodes {
				tag := tags.Eq(i)

				text += fmt.Sprintf("%s ", tag.Text())
			}
			text += fmt.Sprintf("\n")
		}

	}

	return text
}
