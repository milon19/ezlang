package po

import (
	"bufio"
	"fmt"
	"github.com/milon19/ezlang/pkg/translation"
	"os"
	"strings"
)

type CommentType string

const (
	TranslatorComment CommentType = "#"
	ExtractedComment  CommentType = "#."
	ReferenceComment  CommentType = "#:"
	FlagComment       CommentType = "#,"
	PreviousComment   CommentType = "#|"
)

type Comment struct {
	Type    CommentType
	Content string
}

type Entry struct {
	Comments     []Comment
	MsgContext   string
	MsgID        string
	MsgIDPlural  string
	MsgStr       string
	MsgStrPlural []string
	IsFuzzy      bool
	RawLines     []string // To preserve exact formatting
}

type Header struct {
	Comments []Comment
	Metadata map[string]string
	RawLines []string
}

type File struct {
	Header  Header
	Entries []Entry
}

func NewFile() *File {
	return &File{
		Header: Header{
			Metadata: make(map[string]string),
		},
		Entries: make([]Entry, 0),
	}
}

func ReadPoFile(filepath string) (*File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Read all lines first
	var lines []string
	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024 // 1MB buffer
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}

	poFile := NewFile()

	// Phase 1: Parse header
	headerEndIndex := parseHeader(lines, poFile)

	// Phase 2: Parse entries
	parseEntries(lines[headerEndIndex:], poFile)

	return poFile, nil
}

func parseHeader(lines []string, poFile *File) int {
	var headerLines []string
	var i int

	// Collect all lines until we find msgid ""
	for i = 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "msgid \"\"") {
			break
		}
		headerLines = append(headerLines, lines[i])
	}

	// Now collect the header entry (from msgid "" until empty line)
	for ; i < len(lines); i++ {
		headerLines = append(headerLines, lines[i])
		if strings.TrimSpace(lines[i]) == "" {
			break
		}
	}

	poFile.Header.RawLines = headerLines
	processHeaderMetadata(headerLines, poFile)

	// Return the index of the next line after the header
	return i + 1
}

func processHeaderMetadata(headerLines []string, poFile *File) {
	var inMsgStr = false

	for _, line := range headerLines {
		if strings.HasPrefix(line, "#") {
			comment := parseComment(line)
			poFile.Header.Comments = append(poFile.Header.Comments, comment)
		} else if strings.HasPrefix(line, "msgstr") {
			inMsgStr = true
		} else if inMsgStr && strings.HasPrefix(line, "\"") {
			content := extractQuotedText(line)
			if strings.Contains(content, ":") {
				key, value := parseMetadata(content)
				if key != "" {
					poFile.Header.Metadata[key] = value
				}
			}
		}
	}
}

func parseEntries(lines []string, poFile *File) {
	var currentEntry Entry
	var currentField string
	var rawLines []string

	for _, line := range lines {
		// End of entry (empty line)
		if strings.TrimSpace(line) == "" {
			if len(rawLines) > 0 {
				currentEntry.RawLines = rawLines
				poFile.Entries = append(poFile.Entries, currentEntry)
				currentEntry = Entry{}
				rawLines = []string{}
			}
			rawLines = append(rawLines, line)
			continue
		}

		rawLines = append(rawLines, line)

		// Handle comments
		if strings.HasPrefix(line, "#") {
			comment := parseComment(line)
			currentEntry.Comments = append(currentEntry.Comments, comment)
			if comment.Type == FlagComment && strings.Contains(comment.Content, "fuzzy") {
				currentEntry.IsFuzzy = true
			}
			continue
		}

		// Handle message fields
		if strings.HasPrefix(line, "msg") {
			currentField = getMsgFieldType(line)
			updateEntryField(&currentEntry, currentField, extractQuotedText(line))
		} else if strings.HasPrefix(line, "\"") {
			// Continuation line
			appendToEntryField(&currentEntry, currentField, extractQuotedText(line))
		}
	}

	// Add the last entry if exists
	if len(rawLines) > 0 {
		currentEntry.RawLines = rawLines
		poFile.Entries = append(poFile.Entries, currentEntry)
	}
}

func parseComment(line string) Comment {
	if len(line) < 2 {
		return Comment{Type: TranslatorComment, Content: ""}
	}

	var commentType CommentType
	var content string

	if strings.HasPrefix(line, "#.") {
		commentType = ExtractedComment
		content = strings.TrimSpace(line[2:])
	} else if strings.HasPrefix(line, "#:") {
		commentType = ReferenceComment
		content = strings.TrimSpace(line[2:])
	} else if strings.HasPrefix(line, "#,") {
		commentType = FlagComment
		content = strings.TrimSpace(line[2:])
	} else if strings.HasPrefix(line, "#|") {
		commentType = PreviousComment
		content = strings.TrimSpace(line[2:])
	} else {
		commentType = TranslatorComment
		content = strings.TrimSpace(line[1:])
		if content == "" && len(line) == 1 {
			content = ""
		}
	}

	return Comment{
		Type:    commentType,
		Content: content,
	}
}

func parseMetadata(line string) (string, string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// Remove quotes if present
	key = strings.Trim(key, "\"")
	value = strings.Trim(value, "\"")

	return key, value
}

func extractQuotedText(line string) string {
	start := strings.Index(line, "\"")
	end := strings.LastIndex(line, "\"")
	if start >= 0 && end > start {
		return line[start+1 : end]
	}
	return ""
}

func getMsgFieldType(line string) string {
	if strings.HasPrefix(line, "msgctxt") {
		return "msgctxt"
	} else if strings.HasPrefix(line, "msgid_plural") {
		return "msgid_plural"
	} else if strings.HasPrefix(line, "msgid") {
		return "msgid"
	} else if strings.HasPrefix(line, "msgstr[") {
		return "msgstr_plural"
	} else if strings.HasPrefix(line, "msgstr") {
		return "msgstr"
	}
	return ""
}

func updateEntryField(entry *Entry, field string, value string) {
	switch field {
	case "msgctxt":
		entry.MsgContext = value
	case "msgid":
		entry.MsgID = value
	case "msgid_plural":
		entry.MsgIDPlural = value
	case "msgstr":
		entry.MsgStr = value
	case "msgstr_plural":
		if entry.MsgStrPlural == nil {
			entry.MsgStrPlural = make([]string, 0)
		}
		entry.MsgStrPlural = append(entry.MsgStrPlural, value)
	}
}

// Helper function to append to entry fields
func appendToEntryField(entry *Entry, field string, value string) {
	switch field {
	case "msgctxt":
		entry.MsgContext += value
	case "msgid":
		entry.MsgID += value
	case "msgid_plural":
		entry.MsgIDPlural += value
	case "msgstr":
		entry.MsgStr += value
	case "msgstr_plural":
		if len(entry.MsgStrPlural) > 0 {
			lastIdx := len(entry.MsgStrPlural) - 1
			entry.MsgStrPlural[lastIdx] += value
		}
	}
}

func WritePOFile(filepath string, poFile *File, targetLang string, tc *translation.Client) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Printf("error flushing buffer: %v\n", err)
		}
	}(writer)

	// Write header as is
	for _, line := range poFile.Header.RawLines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	// Process and write entries
	for _, entry := range poFile.Entries {
		shouldTranslate := entry.MsgStr == "" || entry.IsFuzzy
		skipNextEmptyLine := false

		for i, line := range entry.RawLines {
			// Skip empty line if previous line was a removed fuzzy flag
			if skipNextEmptyLine && strings.TrimSpace(line) == "" {
				skipNextEmptyLine = false
				continue
			}

			// Handle fuzzy flags
			if strings.HasPrefix(line, "#,") {
				flags := strings.Split(strings.TrimPrefix(line, "#,"), ",")
				var newFlags []string

				for _, flag := range flags {
					flag = strings.TrimSpace(flag)
					if flag != "fuzzy" {
						newFlags = append(newFlags, flag)
					}
				}

				// If there are remaining flags, write them
				if len(newFlags) > 0 {
					_, err := fmt.Fprintf(writer, "#, %s\n", strings.Join(newFlags, ", "))
					if err != nil {
						return err
					}
				} else {
					// If we removed the entire line, skip the next empty line
					skipNextEmptyLine = true
				}
				continue
			}

			// Handle msgstr translation
			if strings.HasPrefix(line, "msgstr \"") && !strings.HasPrefix(line, "msgstr[") {
				if shouldTranslate && entry.MsgID != "" {
					translated := translateText(entry.MsgID, targetLang, tc)
					writeTranslatedString(writer, "msgstr", translated, line)
				} else {
					fmt.Fprintln(writer, line)
				}
				continue
			}

			// Handle plural msgstr translation
			if strings.HasPrefix(line, "msgstr[") {
				if shouldTranslate && (entry.MsgID != "" || entry.MsgIDPlural != "") {
					// Extract the index
					start := strings.Index(line, "[")
					end := strings.Index(line, "]")
					if start != -1 && end != -1 {
						index := line[start : end+1]
						var sourceText string
						if strings.HasPrefix(line, "msgstr[0]") {
							sourceText = entry.MsgID
						} else {
							sourceText = entry.MsgIDPlural
						}
						translated := translateText(sourceText, targetLang, tc)
						writeTranslatedString(writer, "msgstr"+index, translated, line)
					}
				} else {
					fmt.Fprintln(writer, line)
				}
				continue
			}

			// Skip continuation lines for msgstr if we're translating
			if shouldTranslate && strings.HasPrefix(strings.TrimSpace(line), "\"") {
				// Check if this is a continuation of msgstr
				foundMsgStr := false
				for j := i - 1; j >= 0; j-- {
					if strings.HasPrefix(entry.RawLines[j], "msgstr") {
						foundMsgStr = true
						break
					}
					if !strings.HasPrefix(strings.TrimSpace(entry.RawLines[j]), "\"") {
						break
					}
				}
				if foundMsgStr {
					// This is a msgstr continuation, skip it
					continue
				}
			}

			// Write all other lines as is
			_, err := fmt.Fprintln(writer, line)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Helper function to write translated strings (handles multiline)
func writeTranslatedString(writer *bufio.Writer, prefix string, text string, originalLine string) {
	// Get the indentation from the original line
	indent := ""
	for i, char := range originalLine {
		if char != ' ' && char != '\t' {
			indent = originalLine[:i]
			break
		}
	}

	// If the text is short enough, write it on one line
	if len(text) <= 70 {
		fmt.Fprintf(writer, "%s%s \"%s\"\n", indent, prefix, escapeString(text))
		return
	}

	// For longer text, split into multiple lines
	fmt.Fprintf(writer, "%s%s \"\"\n", indent, prefix)

	// Split text into chunks (respecting word boundaries)
	chunks := splitIntoChunks(text, 70)
	for _, chunk := range chunks {
		fmt.Fprintf(writer, "%s\"%s\"\n", indent, escapeString(chunk))
	}
}

// Helper function to split text into chunks
func splitIntoChunks(text string, maxLength int) []string {
	var chunks []string
	words := strings.Fields(text)

	var currentChunk string
	for _, word := range words {
		if len(currentChunk)+len(word)+1 > maxLength && currentChunk != "" {
			chunks = append(chunks, currentChunk)
			currentChunk = word
		} else {
			if currentChunk == "" {
				currentChunk = word
			} else {
				currentChunk += " " + word
			}
		}
	}

	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func translateText(text string, lang string, tc *translation.Client) string {
	fmt.Printf("Translating text: %s\n", text)
	translatedText, err := tc.TranslateText(text, "en", lang)

	if err != nil {
		fmt.Printf("error translating text: %s\n", text)
		return text
	}
	fmt.Printf("Translated text: %s\n", translatedText)
	return translatedText
}
