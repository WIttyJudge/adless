package hostsfile

import (
	"barrier/internal/config"
	"barrier/internal/http"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const localhost = "127.0.0.1"

// Processor is a structure that is responsible for processing blocklists.
type Processor struct {
	config     *config.Config
	httpClient *http.HTTP
}

// Result contains multiple parsed blocklists.
type Result struct {
	startTag           string
	endTag             string
	descriptionComment string
	parsedBlocklists   []ParsedBlocklist
}

// ParsedBlocklist represents a completed result of blocklist
// that is ready to be appended into hosts file.
type ParsedBlocklist struct {
	DomainsCount int

	linesContent []LineContent
}

type LineContent struct {
	ipAddress  string
	domainName string
}

// NewProcessor initializes Processor structure.
func NewProcessor(config *config.Config) *Processor {
	httpClient := http.New()

	return &Processor{
		config:     config,
		httpClient: httpClient,
	}
}

// Process processes blocklists and returns a finished result
// that is ready to save to hosts file.
func (p *Processor) Process() (Result, error) {
	parsedBlocklists := make([]ParsedBlocklist, 0, len(p.config.Blocklists))

	for _, blocklist := range p.config.Blocklists {
		target := blocklist.Target

		log.Info().Str("target", target).Msg("processing blocklist..")

		fileContent, err := p.httpClient.Get(target)
		if err != nil {
			log.Error().Err(err).Str("target", target).Msg("failed to process blocklist")
			continue
		}

		parsedBlocklist := p.processBlocklist(fileContent)
		log.Info().Str("target", target).Msgf("number of %d domains parsed", parsedBlocklist.DomainsCount)

		parsedBlocklists = append(parsedBlocklists, parsedBlocklist)
	}

	result := Result{
		startTag:           StartTag,
		endTag:             EndTag,
		descriptionComment: DescriptionComment,
		parsedBlocklists:   parsedBlocklists,
	}

	return result, nil
}

func (p *Processor) processBlocklist(content string) ParsedBlocklist {
	lines := strings.Split(content, "\n")

	linesContent := make([]LineContent, 0, len(lines))

	for _, line := range lines {
		// remove empty spaces
		line := strings.TrimSpace(line)

		// skip empty lines and comments
		if line == "" || p.isLineComment(line) {
			continue
		}

		line = p.removeInLineComment(line)

		lineContent := LineContent{
			ipAddress: localhost,
		}

		parts := strings.Fields(line)
		if len(parts) == 1 {
			lineContent.domainName = parts[0]
		} else {
			lineContent.domainName = parts[1]
		}

		linesContent = append(linesContent, lineContent)
	}

	parsedBlocklist := ParsedBlocklist{
		linesContent: linesContent,
		DomainsCount: len(linesContent),
	}

	return parsedBlocklist
}

func (p *Processor) isLineComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func (p *Processor) removeInLineComment(line string) string {
	return strings.Split(line, "#")[0]
}

func (r Result) FormatToHostsfile() string {
	var builder strings.Builder

	builder.WriteString(r.startTag)
	builder.WriteString(r.descriptionComment)

	for _, parsedBlocklist := range r.parsedBlocklists {
		for _, lineContent := range parsedBlocklist.linesContent {
			builder.WriteString(lineContent.Format())
		}
	}

	builder.WriteString(r.endTag)

	withoutLastWhitespace := strings.TrimSuffix(builder.String(), "\n")

	return withoutLastWhitespace
}

func (lc LineContent) Format() string {
	return fmt.Sprintf("%s %s\n", lc.ipAddress, lc.domainName)
}
