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
	// *whips out cancer*
	Actions = regexp.MustCompile(`\*.*\*`)
	// Lewdbot
	Lewdbot = regexp.MustCompile(`(?i)lewdbot`)
	// JustPunctuation matches weird junk people send as empty messages.
	JustPunctuation = regexp.MustCompile(`^[\.\\/!?:]`)
	// Kills chatlog pasting
	LeadingNumbers = regexp.MustCompile(`^\d{2,}`)
	// Mentions in discord <@126510493828513793>
	Mentions = regexp.MustCompile(`<@(\d+)>`)
	// Exhentai urls
	GalleryLink = regexp.MustCompile(`http://exhentai.org/g/([[:digit:]]+)/([[:alnum:]]+)`)
	GalleryPage = regexp.MustCompile(`http://exhentai.org/s/([[:alnum:]]+)/([[:digit:]]+)-([[:digit:]]+)`)
)
