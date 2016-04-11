package regex

import "regexp"

var (
	// Russian matches all the cyrillic bullshit they write.
	Russian = regexp.MustCompile(`\p{Cyrillic}`)
	// Link matches inline hyperlinks.
	Link = regexp.MustCompile(`(https?:\/\/[^\s]+)`)
	// Emoticon matches :steamemoticons:. Note the second form of colon that
	// appears if an emoticon is "transformed" into the actual emote.
	Emoticon = regexp.MustCompile(`((:|ː)\w+(:|ː))`)
	// Junk matches... why does this exist, soft?
	Junk = regexp.MustCompile(`[:"]`)
	// WikipediaCitations matches[1] these annoying citation[2] marks.
	WikipediaCitations = regexp.MustCompile(`(\[\d+\])`)
	// RepeatedWhitespace matches 2 or more pieces of whitespace. Make sure to
	// replace them with 1 space instead of nothing!
	RepeatedWhitespace = regexp.MustCompile(`\s{2,}/`)
	// TrailingPunctuation matches any punctuation at the end of the message, to
	// be replaced with a tilde~
	TrailingPunctuation = regexp.MustCompile(`[\.,—\-\~]+$`)
	// NotActualText matches everything that's not Latin text or spaces.
	NotActualText = regexp.MustCompile(`[^\p{L} ]`)
	// Greentext matches '>lines like these'
	Greentext = regexp.MustCompile(`^>`)
	// Actions *whips out cancer*
	Actions = regexp.MustCompile(`\*.*\*`)
	// Lewdbot case insensitive lewdbot match
	Lewdbot = regexp.MustCompile(`(?i)lewdbot`)
	// JustPunctuation matches weird junk people send as empty messages.
	JustPunctuation = regexp.MustCompile(`^[\.\\/!?:]`)
	// LeadingNumbers Kills chatlog pasting
	LeadingNumbers = regexp.MustCompile(`^\d{2,}`)
	// Mentions in discord <@126510493828513793>
	Mentions = regexp.MustCompile(`<@(\d+)>`)
	// ExGalleryLink Matches exhentai gallery links
	ExGalleryLink = regexp.MustCompile(`http:\/\/exhentai\.org\/g\/([[:digit:]]+)/([[:alnum:]]+)`)
	// ExGalleryPage Matches exhentai gallery page links
	ExGalleryPage = regexp.MustCompile(`http:\/\/exhentai\.org\/s\/([[:alnum:]]+)/([[:digit:]]+)-([[:digit:]]+)`)
	// NhGalleryLink Matches exhentai gallery linkss
	NhGalleryLink = regexp.MustCompile(`http:\/\/nhentai\.net\/g\/([[:digit:]]+)`)
)
