package hostsfile

import (
	"barrier/internal/config"
	"barrier/internal/http"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

const localhost = "127.0.0.1"

// Processor is a structure that is responsible for processing blocklists.
type Processor struct {
	wg         *sync.WaitGroup
	config     *config.Config
	httpClient *http.HTTP
}

// Result contains multiple parsed blocklists.
type Result struct {
	startTag           string
	endTag             string
	descriptionComment string
	domainsBlocklist   map[string]LineContent
}

// BlocklistResult represents a parsed result of blocklist
// that is ready to be appended into hosts file.
type BlocklistResult struct {
	DomainsCount int

	linesContent map[string]LineContent
}

type LineContent struct {
	ipAddress  string
	domainName string
}

// NewProcessor initializes Processor structure.
func NewProcessor(config *config.Config) *Processor {
	httpClient := http.New()

	return &Processor{
		wg:         &sync.WaitGroup{},
		config:     config,
		httpClient: httpClient,
	}
}

// Process processes blocklists and returns a finished result
// that is ready to save to hosts file.
func (p *Processor) Process() (Result, error) {
	blocklistsResult := make([]BlocklistResult, len(p.config.Blocklists))
	totalDomainsCount := 0

	for i, blocklist := range p.config.Blocklists {
		i := i
		target := blocklist.Target

		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			blocklistResult, err := p.processBlocklist(target)
			if err != nil {
				log.Error().Err(err).Str("target", target).Msg("failed to process blocklist")
				return
			}

			blocklistsResult[i] = blocklistResult
			totalDomainsCount += blocklistResult.DomainsCount
		}()
	}

	p.wg.Wait()

	domainsBlocklist := p.domainsBlocklist(blocklistsResult)

	result := Result{
		startTag:           StartTag,
		endTag:             EndTag,
		descriptionComment: DescriptionComment,
		domainsBlocklist:   domainsBlocklist,
	}

	log.Info().Msgf("total number of uniq domains: %d", len(domainsBlocklist))

	return result, nil
}

func (p *Processor) processBlocklist(target string) (BlocklistResult, error) {
	log.Info().Str("target", target).Msg("processing blocklist..")

	fileContent, err := p.httpClient.Get(target)
	if err != nil {
		return BlocklistResult{}, err
	}

	linesContent := p.processContent(fileContent)

	blocklistResult := BlocklistResult{
		DomainsCount: len(linesContent),
		linesContent: linesContent,
	}

	log.Info().Str("target", target).Msgf("number of domains: %d", blocklistResult.DomainsCount)

	return blocklistResult, nil
}

func (p *Processor) processContent(content string) map[string]LineContent {
	lines := strings.Split(content, "\n")
	linesContent := make(map[string]LineContent)

	for _, line := range lines {
		// remove empty spaces
		line := strings.TrimSpace(line)

		// skip empty lines, comments, ABP comments and ABP headers
		if line == "" || p.isLineComment(line) || p.isABPComment(line) || p.isABPHeader(line) {
			continue
		}

		// remove a comment in the middle of line
		line = p.removeInLineComment(line)

		// Convert all characters to lowercase
		line = strings.ToLower(line)

		if p.isABPDomain(line) {
			line = p.parseABPDomain(line)
		}

		var domainName string
		parts := strings.Fields(line)
		if len(parts) == 1 {
			domainName = parts[0]
		} else {
			domainName = parts[1]
		}

		lineContent := LineContent{
			ipAddress:  localhost,
			domainName: domainName,
		}

		linesContent[domainName] = lineContent
	}

	return linesContent
}

func (p *Processor) domainsBlocklist(blocklistsResult []BlocklistResult) map[string]LineContent {
	domainsBlocklist := make(map[string]LineContent)

	for _, result := range blocklistsResult {
		for _, line := range result.linesContent {
			domainsBlocklist[line.domainName] = line
		}
	}

	return domainsBlocklist
}

func (p *Processor) isLineComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

// Example: https://raw.githubusercontent.com/hagezi/dns-blocklists/main/adblock/light.txt
func (p *Processor) isABPComment(line string) bool {
	return strings.HasPrefix(line, "!")
}

func (p *Processor) isABPHeader(line string) bool {
	return strings.HasPrefix(line, "[")
}

func (p *Processor) removeInLineComment(line string) string {
	return strings.Split(line, "#")[0]
}

func (p *Processor) isABPDomain(line string) bool {
	return strings.HasPrefix(line, "||") && strings.HasSuffix(line, "^")
}

func (p *Processor) parseABPDomain(line string) string {
	line = strings.TrimPrefix(line, "||")
	line = strings.TrimSuffix(line, "^")
	return line
}

func (r Result) FormatToHostsfile() string {
	var builder strings.Builder

	builder.WriteString(r.startTag)
	builder.WriteString(r.descriptionComment)

	for _, domain := range r.domainsBlocklist {
		builder.WriteString(domain.Format())
	}

	builder.WriteString(r.endTag)

	withoutLastWhitespace := strings.TrimSuffix(builder.String(), "\n")

	return withoutLastWhitespace
}

func (lc LineContent) Format() string {
	return fmt.Sprintf("%s %s\n", lc.ipAddress, lc.domainName)
}
