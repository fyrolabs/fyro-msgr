# Fyro Mailer

## Configuration
```go
type NewMailerOpts struct {
	// Path to email layout, locales, and templates
	TemplatesDir  string
	DefaultFrom   string
	Provider      Provider
	DefaultLocale string
	// Dynamic data to be used in the layout
	LayoutData MailData
}
```

## Letters
Each letter has its own templates and locales.

```go
type RegisterLetterOpts struct {
	Name           string // must be unique and match folder name
	ExtraTemplates []string // index.html.tmpl is included, define extras here
}
```

 Register added letters using:

```go
mailer.RegisterLetter(RegisterLetterOpts{})
```

## Templates
File structure should be as follows:

```
templates
  layout.html.tmpl // Root layout
  locale.en.yml
  locale.zh-cn.yml
  [entryName]/
    locale.en.yml
    locale.zh-cn.yml
    index.html.tmpl
    additional.html.tmpl
```

Locale files should be named locale.[lang].yml

### Subject
Letter locale files must include a mandatory subject message entry, this is templated using data passed in.

```yaml
# YAML
subject: Hello {{ .Name }}
```

## Sending
```go
type SendOpts struct {
	LetterName string
	Data       MailData `json:"data"`
	To         string
	From       string
	ReplyTo    string
	Subject    string // Leave blank to use subject from template
	Locale     string // Leave blank for default
}
```

## Example

Check test/preview.go for a working example
