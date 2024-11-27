package hostsfile

import (
	"adless/internal/config"
	"adless/internal/http"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

const localhost = "127.0.0.1"

// Processor is a structure that is responsible for processing blocklists,
// whitelists and preparing the result to save to hosts file.
type Processor struct {
	config     *config.Config
	httpClient *http.HTTP
}

// Result contains multiple parsed blocklists.
type Result struct {
	startTag           string
	endTag             string
	descriptionComment string
	domains            map[string]LineContent
}

// TargetResult represents a parsed result of blocklist
// that is ready to be appended into hosts file.
type TargetResult struct {
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
		config:     config,
		httpClient: httpClient,
	}
}

// Process processes blocklists and returns a finished result
// that is ready to save to hosts file.
func (p *Processor) Process() (Result, error) {
	wg := &sync.WaitGroup{}
	blocklistsResult := p.processBlocklists(wg)
	whitelistsResult := p.processWhitelists(wg)
	wg.Wait()

	// Merges the results of all targets into one map, where the key
	// is a domain name and the value is a content of the line.
	// Using a domain as a key allows to avoid duplicates.
	blocklistDomains := p.targetDomains(blocklistsResult)
	whitelistDomains := p.targetDomains(whitelistsResult)

	p.applyWhitelist(blocklistDomains, whitelistDomains)

	result := Result{
		startTag:           StartTag,
		endTag:             EndTag,
		descriptionComment: DescriptionComment,
		domains:            blocklistDomains,
	}

	log.Info().Msgf("total number of uniq domains: %d", len(blocklistDomains))

	return result, nil
}

func (p *Processor) processBlocklists(wg *sync.WaitGroup) []TargetResult {
	blocklistsResult := make([]TargetResult, len(p.config.Blocklists))

	for i, blocklist := range p.config.Blocklists {
		i := i
		target := blocklist.Target

		wg.Add(1)

		go func() {
			defer wg.Done()

			blocklistResult, err := p.processBlocklist(target)
			if err != nil {
				log.Error().Err(err).Str("target", target).Msg("failed to process blocklist")
				return
			}

			blocklistsResult[i] = blocklistResult
		}()
	}

	return blocklistsResult
}

func (p *Processor) processWhitelists(wg *sync.WaitGroup) []TargetResult {
	whitelistsResult := make([]TargetResult, len(p.config.Blocklists))

	for i, whitelist := range p.config.Whitelists {
		i := i
		target := whitelist.Target

		wg.Add(1)

		go func() {
			defer wg.Done()

			whitelistResult, err := p.processWhitelist(target)
			if err != nil {
				log.Error().Err(err).Str("target", target).Msg("failed to process whitelist")
				return
			}

			whitelistsResult[i] = whitelistResult
		}()
	}

	return whitelistsResult
}

func (p *Processor) processBlocklist(target string) (TargetResult, error) {
	log.Info().Str("target", target).Msg("processing blocklist..")

	blocklistResult, err := p.proccessListTarget(target)
	if err != nil {
		return TargetResult{}, err
	}

	log.Info().Str("target", target).Msgf("number of domains: %d", blocklistResult.DomainsCount)

	return blocklistResult, nil
}

func (p *Processor) processWhitelist(target string) (TargetResult, error) {
	log.Info().Str("target", target).Msg("processing whitelist..")

	whitelistResult, err := p.proccessListTarget(target)
	if err != nil {
		return TargetResult{}, err
	}

	log.Info().Str("target", target).Msgf("number of domains: %d", whitelistResult.DomainsCount)

	return whitelistResult, nil
}

func (p *Processor) proccessListTarget(target string) (TargetResult, error) {
	fileContent, err := p.httpClient.Get(target)
	if err != nil {
		return TargetResult{}, err
	}

	linesContent := p.processContent(fileContent)

	blocklistResult := TargetResult{
		DomainsCount: len(linesContent),
		linesContent: linesContent,
	}

	return blocklistResult, nil
}

func (p *Processor) processContent(content string) map[string]LineContent {
	lines := strings.Split(content, "\n")
	linesContent := make(map[string]LineContent)

	for _, rawLine := range lines {
		line := p.normalizeLine(rawLine)

		if p.shouldSkipLine(line) {
			continue
		}

		domainName := p.extractDomain(line)
		if p.IsSkippedDomain(domainName) {
			continue
		}

		lineContent := LineContent{
			ipAddress:  localhost,
			domainName: domainName,
		}

		linesContent[domainName] = lineContent
	}

	return linesContent
}

func (p *Processor) targetDomains(targetResult []TargetResult) map[string]LineContent {
	targetDomains := make(map[string]LineContent)

	for _, result := range targetResult {
		for _, line := range result.linesContent {
			targetDomains[line.domainName] = line
		}
	}

	return targetDomains
}

func (p *Processor) applyWhitelist(blocklistDomains, whitelistDomains map[string]LineContent) {
	for key := range whitelistDomains {
		delete(blocklistDomains, key)
	}
}

// 1. Remove empty spaces.
// 2. Remove a comment in the middle of line.
// 3. Convert all characters to lowercase.
func (p *Processor) normalizeLine(rawLine string) string {
	line := strings.TrimSpace(rawLine)
	line = p.removeInLineComment(line)
	line = strings.ToLower(line)

	return line
}

// skip empty lines, comments, ABP comments and ABP headers.
func (p *Processor) shouldSkipLine(line string) bool {
	return line == "" || p.isLineComment(line) || p.isABPComment(line) || p.isABPHeader(line)
}

func (p *Processor) removeInLineComment(line string) string {
	return strings.Split(line, "#")[0]
}

func (p *Processor) parseABPDomain(line string) string {
	line = strings.TrimPrefix(line, "||")
	line = strings.TrimSuffix(line, "^")
	return line
}

func (p *Processor) extractDomain(line string) string {
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

	return domainName
}

func (p *Processor) isLineComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

// supported ABP style: ||subdomain.domain.tlp^
// Example: https://raw.githubusercontent.com/hagezi/dns-blocklists/main/adblock/light.txt
func (p *Processor) isABPDomain(line string) bool {
	return strings.HasPrefix(line, "||") && strings.HasSuffix(line, "^")
}

func (p *Processor) isABPComment(line string) bool {
	return strings.HasPrefix(line, "!")
}

func (p *Processor) isABPHeader(line string) bool {
	return strings.HasPrefix(line, "[")
}

// IsSkippedDomain checks if a domain is in the skip list.
// Some lists (i.e StevenBlack's) contain these as they are supposed to be used as HOST.
func (p *Processor) IsSkippedDomain(domain string) bool {
	skipList := []string{
		"localhost",
		"localhost.localdomain",
		"local",
		"broadcasthost",
		"ip6-localhost",
		"ip6-loopback",
		"lo0 localhost",
		"ip6-localnet",
		"ip6-mcastprefix",
		"ip6-allnodes",
		"ip6-allrouters",
		"ip6-allhosts",
	}

	return slices.Contains(skipList, domain)
}

func (r Result) FormatToHostsfile() string {
	var builder strings.Builder

	builder.WriteString(r.startTag)
	builder.WriteString(r.descriptionComment)

	for _, domain := range r.domains {
		builder.WriteString(domain.Format())
	}

	builder.WriteString(r.endTag)

	withoutLastWhitespace := strings.TrimSuffix(builder.String(), "\n")

	return withoutLastWhitespace
}

func (lc LineContent) Format() string {
	return fmt.Sprintf("%s %s\n", lc.ipAddress, lc.domainName)
}
