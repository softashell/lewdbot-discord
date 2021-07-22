package mangadex

import "time"

type Manga struct {
	Result string `json:"result"`
	Data   struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title map[string]string `json:"title"`
			Description map[string]string `json:"description"`
			OriginalLanguage       string `json:"originalLanguage"`
			LastVolume             string `json:"lastVolume"`
			LastChapter            string `json:"lastChapter"`
			PublicationDemographic string `json:"publicationDemographic"`
			Status                 string `json:"status"`
			Year                   int    `json:"year"`
			ContentRating          string `json:"contentRating"`
			Tags                   []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					Name map[string]string `json:"name"`
					Group   string `json:"group"`
					Version int    `json:"version"`
				} `json:"attributes"`
			} `json:"tags"`
			Version   int    `json:"version"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
		} `json:"attributes"`

	} `json:"data"`
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
	Errors []struct {
		Detail string `json:"detail"`
		ID     string `json:"id"`
		Status int    `json:"status"`
		Title  string `json:"title"`
	} `json:"errors"`
}

type Chapter struct {
	Data struct {
		Attributes struct {
			Chapter            string    `json:"chapter"`
			CreatedAt          time.Time `json:"createdAt"`
			Data               []string  `json:"data"`
			DataSaver          []string  `json:"dataSaver"`
			Hash               string    `json:"hash"`
			PublishAt          time.Time `json:"publishAt"`
			Title              string    `json:"title"`
			TranslatedLanguage string    `json:"translatedLanguage"`
			UpdatedAt          time.Time `json:"updatedAt"`
			Version            int       `json:"version"`
			Volume             string    `json:"volume"`
		} `json:"attributes"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
	Errors []struct {
		Detail string `json:"detail"`
		ID     string `json:"id"`
		Status int    `json:"status"`
		Title  string `json:"title"`
	} `json:"errors"`
	Result string `json:"result"`
}

type Cover struct {
	Data struct {
		Attributes struct {
			CreatedAt   string `json:"createdAt"`
			Description string `json:"description"`
			FileName    string `json:"fileName"`
			UpdatedAt   string `json:"updatedAt"`
			Version     int    `json:"version"`
			Volume      string `json:"volume"`
		} `json:"attributes"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Relationships []struct {
		Attributes struct {
		} `json:"attributes"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
	Errors []struct {
		Detail string `json:"detail"`
		ID     string `json:"id"`
		Status int    `json:"status"`
		Title  string `json:"title"`
	} `json:"errors"`
	Result string `json:"result"`
}

type Author struct {
	Data struct {
		Attributes struct {
			CreatedAt string `json:"createdAt"`
			ImageURL  string `json:"imageUrl"`
			Name      string `json:"name"`
			UpdatedAt string `json:"updatedAt"`
			Version   int    `json:"version"`
		} `json:"attributes"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Relationships []struct {
		Attributes struct {
		} `json:"attributes"`
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
	Result string `json:"result"`
}