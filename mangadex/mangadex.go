package mangadex

import (
	"errors"
	"fmt"
	"html"
	"time"

	"github.com/bwmarrin/discordgo"
	lru "github.com/hashicorp/golang-lru"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
	"github.com/softashell/lewdbot-discord/regex"
)

const apiUrl = "https://api.mangadex.org"

type Mangadex struct{
	mangaCache *lru.TwoQueueCache
	chapterCache *lru.TwoQueueCache
	coverCache *lru.TwoQueueCache
	authorCache *lru.TwoQueueCache
}

func NewMangadex() *Mangadex {
	mangaCache, err := lru.New2Q(256)
	if(err != nil) {
		log.Fatal(err)
	}

	chapterCache, err := lru.New2Q(256)
	if(err != nil) {
		log.Fatal(err)
	}

	coverCache, err := lru.New2Q(256)
	if(err != nil) {
		log.Fatal(err)
	}

	authorCache, err := lru.New2Q(256)
	if(err != nil) {
		log.Fatal(err)
	}

	md := &Mangadex{
		mangaCache: mangaCache,
		chapterCache: chapterCache,
		coverCache: coverCache,
		authorCache: authorCache,
	}

	return md
}

func (md *Mangadex)GetManga(uuid string) Manga {
	var response Manga

	if val, ok := md.mangaCache.Get(uuid); ok {
		response = val.(Manga)

		return response
	}

	requestUrl := apiUrl + "/manga/" + uuid

	_, reply, errs := gorequest.New().Get(requestUrl).Timeout(10 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)
	}

	md.mangaCache.Add(uuid, response)

	return response
}

func (md *Mangadex)GetCover(uuid string) Cover {
	var response Cover

	if val, ok := md.coverCache.Get(uuid); ok {
		response = val.(Cover)

		return response
	}

	requestUrl := apiUrl + "/cover/" + uuid

	_, reply, errs := gorequest.New().Get(requestUrl).Timeout(10 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)
	}

	md.coverCache.Add(uuid, response)

	return response
}

func (md *Mangadex)GetAuthor(uuid string) Author {
	var response Author

	if val, ok := md.authorCache.Get(uuid); ok {
		response = val.(Author)

		return response
	}

	requestUrl := apiUrl + "/author/" + uuid

	_, reply, errs := gorequest.New().Get(requestUrl).Timeout(10 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)
	}

	md.authorCache.Add(uuid, response)

	return response
}

func (md *Mangadex)GetMangaCoverArt(uuid string) string {
	var response string

	m := md.GetManga(uuid)

	for _, rel := range m.Data.Relationships {
		if(rel.Type != "cover_art") {
			continue
		}

		cover := md.GetCover(rel.ID)
		

		response = "https://uploads.mangadex.org/covers/" + m.Data.ID + "/" + cover.Data.Attributes.FileName
		break;
	}

	return response
}

func (md *Mangadex)GetMangaAuthor(uuid string) string {
	var response string

	m := md.GetManga(uuid)

	for _, rel := range m.Data.Relationships {
		if(rel.Type != "author") {
			continue
		}

		author := md.GetAuthor(rel.ID)
		response = author.Data.Attributes.Name
		break;
	}

	return response
}


func (md *Mangadex)GetChapter(uuid string) (Chapter, error) {
	var response Chapter

	if val, ok := md.chapterCache.Get(uuid); ok {
		response = val.(Chapter)

		return response, nil
	}

	requestUrl := apiUrl + "/chapter/" + uuid

	log.Info(requestUrl)

	_, reply, errs := gorequest.New().Get(requestUrl).Timeout(10 * time.Second).EndStruct(&response)
	for _, err := range errs {
		log.WithFields(log.Fields{
			"reply": string(reply),
		}).Error("API Request failed", err)
	}

	md.chapterCache.Add(uuid, response)

	return response, nil
}

func (md *Mangadex)GetChapterManga(uuid string) (Manga, error) {
	var response Manga

	c, err := md.GetChapter(uuid)
	if(err != nil) {
		return response, err
	}

	var mangaUuid string

	for _, rel := range c.Data.Relationships {
		if(rel.Type != "manga") {
			continue
		}

		mangaUuid = rel.ID
		break;
	}

	if(mangaUuid == "") {
		return response, errors.New("No manga ID found")
	}

	response = md.GetManga(mangaUuid)
	return response, nil
}

// ParseLinks Returns gallery metadata from founds links in input
func (md *Mangadex) ParseLinks(s *discordgo.Session, m *discordgo.MessageCreate, channel string, text string) bool {
	var mangaList []Manga
	var chapterList []Chapter

	mdTitleLinks := regex.MangadexTitleUrl.FindAllStringSubmatch(text, -1)
	mdChapterLinks := regex.MangadexChapterUrl.FindAllStringSubmatch(text, -1)

	for _, link := range mdTitleLinks {
		id := link[1]

		manga := md.GetManga(id)

		mangaList = append(mangaList, manga)
	}

	for _, link := range mdChapterLinks {
		id := link[1]

		chapter, err := md.GetChapter(id)
		if(err != nil) {
			continue
		}

		manga, err := md.GetChapterManga(id)
		if(err != nil) {
			continue
		}

		var dupe bool
		for _, v := range mangaList {
			if v.Data.ID == manga.Data.ID {
				dupe = true
				break
			}
		}
		if(!dupe) {
			mangaList = append(mangaList, manga)
		}

		chapterList = append(chapterList, chapter)
	}

	if len(mangaList) > 0 {
		md.writeMangaData(s, m, mangaList)
	} else if len(chapterList) > 0 {
		// Nothing for now, maybe later I will add chapter number to title
	} else {
		// Didn't find anything
		return false
	}

	return true
}


func (md *Mangadex) writeMangaData(s *discordgo.Session, m *discordgo.MessageCreate, mangaList []Manga) {
	for _, manga := range mangaList {
		if len(manga.Errors) > 0 {
			log.Warnf("Manga: %s error: %q", manga.Data.ID, manga.Errors)
		}

		var tags []string

		for _, _tag := range manga.Data.Attributes.Tags {
			if val, ok := _tag.Attributes.Name["en"]; ok {
				tags = append(tags, val)
			}
		}

		var fields []*discordgo.MessageEmbedField

		fields = append(fields, &discordgo.MessageEmbedField{Name: "Author", Value: md.GetMangaAuthor(manga.Data.ID), Inline: true})

		var text string
		for i, tag := range tags {
			if i < len(tags)-1 {
				text += fmt.Sprintf("%s, ", tag)
			} else {
				text += fmt.Sprintf("%s", tag)
			}
		}

		fields = append(fields, &discordgo.MessageEmbedField{Name: "Tags", Value: text, Inline: true})

		var title string
		if val, ok := manga.Data.Attributes.Title["en"]; ok {
			title = val
		}

		var description string
		if val, ok := manga.Data.Attributes.Description["en"]; ok {
			description = val
		}

		message := discordgo.MessageEmbed{
			URL:       fmt.Sprintf("https://mangadex.org/title/%s/", manga.Data.ID),
			Title:     html.UnescapeString(title),
			Description: html.UnescapeString(description),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: md.GetMangaCoverArt(manga.Data.ID)},
			Fields:    fields,
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &message)
		if err != nil {
			log.Warn("s.ChannelMessageSendEmbed >>", err)
			return
		}

		original, err := s.ChannelMessage(m.ChannelID, m.ID)
		if err != nil {
			log.Warn("s.ChannelMessage >>", err)
			return
		}
		edit := discordgo.NewMessageEdit(m.ChannelID, m.ID)
		edit.Flags = original.Flags | discordgo.MessageFlagsSuppressEmbeds
		_, err = s.ChannelMessageEditComplex(edit)
		if err != nil {
			log.Warn("s.ChannelMessageEditComplex >>", err)
			return
		}
	}
}

