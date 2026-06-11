package library

// Common.yaml dosyasından parse edilen yapı taslakları.
// Bu dosya yalnızca tip tanımları içerir; hiçbir iş mantığı barındırmaz.

// Embed üst satırı şablonu: ikon + etiket.
type EmbedAuthor struct {
	Name    string `yaml:"name"`
	IconURL string `yaml:"icon_url"`
}

// Embed veri alanı şablonu.
type EmbedField struct {
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Inline bool   `yaml:"inline"`
}

// Embed alt bilgi satırı şablonu.
type EmbedFooter struct {
	Text    string `yaml:"text"`
	IconURL string `yaml:"icon_url"`
}

// Tek bir embed şablonunun tam yapısı.
// Common.yaml > templates ve commands altındaki her entry bu yapıya parse edilir.
type EmbedTemplate struct {
	Color       string       `yaml:"color"`
	Author      EmbedAuthor  `yaml:"author"`
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Fields      []EmbedField `yaml:"fields"`
	Footer      EmbedFooter  `yaml:"footer"`
}

// Bir zaman biriminin tam ve kısa adı.
type DurationUnit struct {
	Full  string `yaml:"full"`
	Short string `yaml:"short"`
}

// Common.yaml dosyasının tamamını temsil eder.
type CommonLib struct {
	Icons     map[string]string        `yaml:"icons"`
	Templates map[string]EmbedTemplate `yaml:"templates"`
	Commands  map[string]EmbedTemplate `yaml:"commands"`
	Titles    map[string]string        `yaml:"titles"`
	Messages  map[string]string        `yaml:"messages"`
	Duration  map[string]DurationUnit  `yaml:"duration"`
}
