package translation

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"strings"
)

type Client struct {
	client *translate.Client
}

func NewTranslationClient(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	_, err = cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("AWS credentials not found. Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables: %w", err)
	}

	client := translate.NewFromConfig(cfg)
	return &Client{client: client}, nil
}

// TranslateText translates text using AWS Translate
func (tc *Client) TranslateText(text string, sourceLang, targetLang string) (string, error) {
	// AWS Translate has a limit of 10,000 bytes per request
	// If text is too long, split it into chunks
	if len(text) > 9000 {
		return tc.translateLongText(text, sourceLang, targetLang)
	}

	input := &translate.TranslateTextInput{
		Text:               aws.String(text),
		SourceLanguageCode: aws.String(sourceLang),
		TargetLanguageCode: aws.String(targetLang),
	}

	result, err := tc.client.TranslateText(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("translation failed: %w", err)
	}

	return *result.TranslatedText, nil
}

// translateLongText handles texts longer than AWS limit by splitting
func (tc *Client) translateLongText(text string, sourceLang, targetLang string) (string, error) {
	// Split text by paragraphs or sentences
	paragraphs := strings.Split(text, "\n\n")
	var translatedParts []string

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			translatedParts = append(translatedParts, "")
			continue
		}

		// If paragraph is still too long, split by sentences
		if len(paragraph) > 9000 {
			sentences := splitIntoSentences(paragraph)
			var translatedSentences []string

			for _, sentence := range sentences {
				translated, err := tc.TranslateText(sentence, sourceLang, targetLang)
				if err != nil {
					return "", err
				}
				translatedSentences = append(translatedSentences, translated)
			}

			translatedParts = append(translatedParts, strings.Join(translatedSentences, " "))
		} else {
			translated, err := tc.TranslateText(paragraph, sourceLang, targetLang)
			if err != nil {
				return "", err
			}
			translatedParts = append(translatedParts, translated)
		}
	}

	return strings.Join(translatedParts, "\n\n"), nil
}

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
	// Simple sentence splitting (you might want to improve this)
	sentences := strings.Split(text, ". ")
	var result []string

	for i, sentence := range sentences {
		if i < len(sentences)-1 {
			result = append(result, sentence+".")
		} else {
			result = append(result, sentence)
		}
	}

	return result
}
