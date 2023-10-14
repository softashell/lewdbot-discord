package mangadex

import (
	"reflect"
	"testing"

	"github.com/softashell/lewdbot-discord/mangadex"
)

func TestMangadex_GetMangaTitle(t *testing.T) {
	tests := []struct {
		name string
		uuid string
		want string
	}{
		{
			name: "The Bride of Barbaroi",
			uuid: "178d3cd8-4624-4afc-81b0-b2f260a1b28b",
			want: "Hime Kishi wa Barbaroi no Yome",
		},
	}
	m := mangadex.NewMangadex()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.GetManga(tt.uuid)

			var title string
			if val, ok := got.Data.Attributes.Title["en"]; ok {
				title = val
			}
			if !reflect.DeepEqual(title, tt.want) {
				t.Errorf("Mangadex.GetManga() = %q, want %q", title, tt.want)
			}
		})
	}
}

func TestMangadex_GetChapterMangaTitle(t *testing.T) {
	tests := []struct {
		name string
		uuid string
		want string
	}{
		{
			name: "The Bride of Barbaroi",
			uuid: "d2931720-f9f0-4f58-ae88-76e852de61e4",
			want: "Hime Kishi wa Barbaroi no Yome",
		},
	}
	m := mangadex.NewMangadex()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.GetChapterManga(tt.uuid)
			if err != nil {
				t.Error(err)
			}

			var title string
			if val, ok := got.Data.Attributes.Title["en"]; ok {
				title = val
			}
			if !reflect.DeepEqual(title, tt.want) {
				t.Errorf("Mangadex.GetManga() = %q, want %q", title, tt.want)
			}
		})
	}
}
