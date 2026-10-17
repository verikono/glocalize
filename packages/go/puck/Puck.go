package puck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	GT "github.com/verikono/glocalize/types/go"
)

type BytesSize = int

type Puck struct {
	url      GT.Url
	glob     GlobPattern
	protocol GT.UrlProtocol
	format   GT.DataFormat
	resource ResourceType
	IPuck
}

func NewPuck(options ...GT.OptionerFn) *Puck {
	puck := Puck{
		url:      "file:///dev/null",
		format:   FormatUnknown,
		protocol: ProtocolFile,
		resource: UnknownResource,
	}

	for _, option := range options {
		err := option(&puck)
		if err != nil {
			panic(err)
		}
	}

	if err := puck.assertValidPuck(); err != nil {
		panic(err)
	}

	return &puck
}

func (p *Puck) Sync() error {
	return p.assertValidPuck()
}

func (p *Puck) Read() (UTF8FileContent, error) {
	if err := p.assertValidPuck(); err != nil {
		return "", err
	}

	if p.glob != "" {
		return p.readGlob()
	}

	switch p.resolveResource() {
	case FileResource:
		bytes, err := os.ReadFile(p.Path())
		if err != nil {
			return "", fmt.Errorf("failed reading %s: %w", p.Path(), err)
		}
		return UTF8FileContent(bytes), nil
	case DirectoryResource:
		return "", fmt.Errorf("puck cannot read: %s is a directory", p.Path())
	default:
		return "", fmt.Errorf("puck cannot read: unsupported resource type %s", p.resource)
	}
}

func (p *Puck) readGlob() (UTF8FileContent, error) {
	matches, err := filepath.Glob(p.glob)
	if err != nil {
		return "", fmt.Errorf("invalid glob %s: %w", p.glob, err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("glob %s matched no files", p.glob)
	}

	sort.Strings(matches)

	var content strings.Builder
	for _, match := range matches {
		bytes, err := os.ReadFile(match)
		if err != nil {
			return "", fmt.Errorf("failed reading %s: %w", match, err)
		}
		content.Write(bytes)
		content.WriteString("\n")
	}

	return UTF8FileContent(content.String()), nil
}

// Deprecated: with glob's being sensed Puck::ReadGlob can be accomplished via the Puck::Read function
func (p *Puck) ReadGlob(glob GlobPattern) (UTF8FileContent, error) {
	return "", fmt.Errorf("called ReadGlob --- deprecated: use Read")
}

func (p *Puck) Exists() bool {
	switch p.protocol {
	case ProtocolFile:
		{
			_, err := os.Stat(p.Path())
			if err != nil {
				if os.IsNotExist(err) {
					return false
				}
				panic(fmt.Errorf("file system error - %s", err))
			}
			return true // it didn't error which means it has information about the file, thus it exists.
		}
	default:
		{
			panic(fmt.Errorf("cannot calll exists on protocol %s", p.protocol))
		}
	}
}

func (p *Puck) Extension() GT.FileExtension {
	locator := p.Path()
	if p.glob != "" {
		locator = p.glob
	}

	if p.protocol == ProtocolFile && p.resolveResource() == DirectoryResource {
		return ""
	}

	ext := filepath.Ext(filepath.Base(locator))
	return GT.FileExtension(strings.TrimPrefix(ext, "."))
}

func (p *Puck) Path() AbsoluteFsPath {
	switch p.protocol {
	case ProtocolFile:
		if p.glob != "" {
			return ""
		}
		return p.url[7:]
	default:
		return ""
	}
}

func (p *Puck) Url() Url {
	return p.url
}

func (p *Puck) Protocol() UrlProtocol {
	if len(p.url) > 0 {
		parts := strings.Split(p.url, "://")
		if GT.IsUrlProtocol(parts[0]) {
			return UrlProtocol(parts[0])
		}
	}
	panic(fmt.Errorf("invalid url: %s", p.url))
}

func (p *Puck) Format() DataFormat {
	// Agent: when glob is truthy then the glob must end with a shared suffix
	//        indicating the glob pulls back a common format
	if len(p.url) > 0 {
		parts := strings.Split(p.url, "://")
		basepath := parts[len(parts)-1]
		segs := strings.Split(basepath, ".")
		if len(segs) < 2 {
			p.format = FormatUnknown
			return p.format
		}
		ext := segs[len(segs)-1]
		if format, ok, _ := GT.IsDataFormat(ext); ok {
			p.format = format
		} else {
			p.format = FormatUnknown
		}
		return p.format
	}
	panic(fmt.Errorf("invalid url: %s", p.url))
}

func (p *Puck) assertValidPuck() error {
	if err := p.assertValidGlob(); err != nil {
		return err
	}
	if err := p.assertValidPath(); err != nil {
		return err
	}
	if err := p.assertValidProtocol(); err != nil {
		return err
	}
	if err := p.assertValidFormat(); err != nil {
		return err
	}
	if err := p.assertValidUrl(); err != nil {
		return err
	}
	return nil
}

func (p *Puck) assertValidGlob() error {
	if p.glob == "" {
		return nil
	}
	if filepath.Ext(p.glob) == "" {
		return fmt.Errorf("globs must contain only files of the same format and file suffix: %s", p.glob)
	}
	return nil
}

func (p *Puck) assertValidPath() error {
	if p.protocol != ProtocolFile {
		return nil
	}
	if p.glob != "" {
		return nil
	}
	if !filepath.IsAbs(p.Path()) {
		return fmt.Errorf("path must be an absolute file system path: %s", p.Path())
	}
	return nil
}

func (p *Puck) assertValidUrl() error {
	if p.protocol != ProtocolFile {
		return nil
	}
	locator := p.Path()
	if p.glob != "" {
		locator = p.glob
	}
	if p.url != fmt.Sprintf("%s://%s", p.protocol, locator) {
		return fmt.Errorf("url %s does not match %s://%s", p.url, p.protocol, locator)
	}
	return nil
}

func (p *Puck) assertValidProtocol() error {
	if !GT.IsUrlProtocol(string(p.Protocol())) {
		return fmt.Errorf("unsupported protocol: %s", p.Protocol())
	}
	return nil
}

func (p *Puck) assertValidFormat() error {
	if p.protocol != ProtocolFile {
		return nil
	}

	switch p.resolveResource() {
	case DirectoryResource:
		return nil
	case FileResource:
		if p.Format() == FormatUnknown {
			return fmt.Errorf("unknown format for url: %s", p.url)
		}
		return nil
	default:
		return nil
	}
}

func (p *Puck) ResourceType() ResourceType {
	p.resource = p.resolveResource()
	return p.resource
}

// Size returns the total size in bytes of the source: a file's size, the
// summed sizes of every file matched by a glob, or the summed sizes of every
// file beneath a directory.
func (p *Puck) Size() BytesSize {
	if p.protocol != ProtocolFile {
		panic(fmt.Errorf("Puck::Size: unsupported protocol %s", p.protocol))
	}

	if p.glob != "" {
		matches, err := filepath.Glob(p.glob)
		if err != nil {
			panic(fmt.Errorf("Puck::Size: invalid glob %s: %w", p.glob, err))
		}
		var total BytesSize
		for _, match := range matches {
			total += p.sizeOf(match)
		}
		return total
	}

	return p.sizeOf(p.Path())
}

func (p *Puck) sizeOf(path string) BytesSize {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(fmt.Errorf("Puck::Size: failed stat %s: %w", path, err))
	}

	if !fileInfo.IsDir() {
		return BytesSize(fileInfo.Size())
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(fmt.Errorf("Puck::Size: failed reading directory %s: %w", path, err))
	}

	var total BytesSize
	for _, entry := range entries {
		total += p.sizeOf(filepath.Join(path, entry.Name()))
	}
	return total
}

func (p *Puck) resolveResource() ResourceType {
	if p.glob != "" {
		return GlobFileResource
	}

	switch p.protocol {
	case ProtocolFile:
		fileInfo, err := os.Stat(p.Path())
		if err != nil {
			if os.IsNotExist(err) {
				return UnknownResource
			}
			panic(fmt.Errorf("failed obtaining resource type: %s", err))
		}
		if fileInfo.IsDir() {
			return DirectoryResource
		}
		return FileResource
	default:
		panic(fmt.Errorf("failed obtaining resource type: unsupported protocol %s", p.protocol))
	}
}
