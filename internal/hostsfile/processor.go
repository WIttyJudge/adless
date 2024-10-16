package hostsfile

import (
	"barrier/internal/config"
	"barrier/internal/http"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

const localhost = "127.0.0.1"

type Processor struct {
	config     *config.Config
	httpClient *http.HTTP
}

type Result struct {
	ParsedBlocklists []ParsedBlocklist
}

type ParsedBlocklist struct {
	LinesContent []LineContent
}

type LineContent struct {
	ipAddress  string
	domainName string
}

func NewProcessor(config *config.Config) *Processor {
	httpClient := http.New(&config.HTTP)

	return &Processor{
		config:     config,
		httpClient: httpClient,
	}
}

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
		parsedBlocklists = append(parsedBlocklists, parsedBlocklist)
	}

	result := Result{
		ParsedBlocklists: parsedBlocklists,
	}

	return result, nil
}

func (p *Processor) processBlocklist(content string) ParsedBlocklist {
	lines := strings.Split(content, "\n")

	linesContent := make([]LineContent, 0, len(lines))

	for _, line := range lines {
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
		LinesContent: linesContent,
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

	for _, parsedBlocklist := range r.ParsedBlocklists {
		for _, lineContent := range parsedBlocklist.LinesContent {
			builder.WriteString(lineContent.Format())
		}
	}

	withoutLastWhitespace := strings.TrimSuffix(builder.String(), "\n")

	return withoutLastWhitespace
}

func (lc LineContent) Format() string {
	return fmt.Sprintf("%s %s\n", lc.ipAddress, lc.domainName)
}
