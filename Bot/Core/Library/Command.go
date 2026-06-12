package library

// Komut şablonu.
type CommandLib struct {
	Name          string            `yaml:"name"`
	Description   string            `yaml:"description"`
	Usage         string            `yaml:"usage"`
	Examples      map[string]string `yaml:"examples"`
	Titles        map[string]string `yaml:"titles"`
	Messages      map[string]string `yaml:"messages"`
	Icons         map[string]string `yaml:"icons"`
	Options       CommandOptions    `yaml:"options"`
	Fields        []CommandField    `yaml:"fields,omitempty"`
	Footers       map[string]string `yaml:"footers,omitempty"`
	IsCoreCommand bool              // YAML dan okunmaz , komutlar oluşturulurken doldurulur. Sistem komutları ile Modül komutlarını ayırmak için kullanılır.
}

// Komut şablonu içerisindeki Fields alanı şablonu.
type CommandField struct {
	ID     int    `yaml:"id"`
	Name   string `yaml:"name"`
	Value  string `yaml:"value"`
	Inline bool   `yaml:"inline"`
}

// Komut şablonu içerisindeki Options alanı şablonu.
type CommandOptions struct {
	// Genel — her komutta ortak
	Enabled   bool   `yaml:"enabled"`
	Category  string `yaml:"category"`
	Ephemeral bool   `yaml:"ephemeral"`
	Cooldown  int    `yaml:"cooldown"`

	// Komuta özel scalar değerler — GetOptionInt/String/Bool buradan okur
	Extra map[string]any `yaml:",inline"`

	// Typed listeler — her zaman aynı yapıda
	AutoActions []AutoAction `yaml:"auto_actions"`
	Slash       []SlashArg   `yaml:"slash"`
}

// Otomatik aksiyon alanı şablonu.
type AutoAction struct {
	At     int    `yaml:"at"`
	Action string `yaml:"action"`
	Reason string `yaml:"reason"`
}

// Slash komut kullanımı şablonu.
type SlashArg struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Type        string         `yaml:"type"`
	Required    bool           `yaml:"required"`
	MinValue    int            `yaml:"min_value"`
	MaxValue    int            `yaml:"max_value"`
	Extra       map[string]any `yaml:",inline"`
}

// Komut adını döndürür.
func (c *CommandLib) GetName() string { return c.Name }

// Komut açıklamasını döndürür.
func (c *CommandLib) GetDescription() string { return c.Description }

// Komut kullanım şeklini döndürür.
func (c *CommandLib) GetUsage() string { return c.Usage }

// Komut field alanını döndürür
func (c *CommandLib) GetField(tagID int) *CommandField {

	for _, v := range c.Fields {

		if v.ID == tagID {
			return &v
		}
	}
	return nil
}

// Komut başlıklarından verilen anahtardaki başlığı döndürür.
func (c *CommandLib) GetTitle(key, fallback string) string {
	if v, ok := c.Titles[key]; ok {
		return v
	}
	return fallback
}

// Komut mesajlarından verilen anahtardaki mesajı döndürür.
func (c *CommandLib) GetMessage(key, fallback string) string {
	if v, ok := c.Messages[key]; ok {
		return v
	}
	return fallback
}

// Komut ikon URL adreslerinden verilen anahtardaki URL adresini döndürür.
func (c *CommandLib) GetIcon(key, fallback string) string {
	if v, ok := c.Icons[key]; ok {
		return v
	}
	return fallback
}

// Komut Seçeneklerini tip dönüştürür.
func GetOptionsExtraMap[T any](c *CommandLib, key string) (map[string]T, bool) {
	val, ok := c.Options.Extra[key]

	m, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}
	result := make(map[string]T)
	for k, v := range m {
		if typedVal, ok := v.(T); ok {
			result[k] = typedVal
		}
	}
	return result, true
}

// Ekstra alanlardan veri çekerken tip dönüşümlerini yönetir
func getExtra[T any](c *CommandLib, key string) (val T, ok bool) {
	v, _ok := c.Options.Extra[key]
	if !_ok {
		var zero T
		return zero, false
	}

	// Eğer beklenen tip int ise ve gelen float64 ise dönüştür
	if val, _ok := v.(T); _ok {
		return val, true
	}

	// YAML'dan gelen sayısal veri float64 ise ve T int ise özel kontrol
	if f, _ok := v.(float64); _ok {
		var dummy any = *new(T)
		if _, _ok := dummy.(int); _ok {
			return any(int(f)).(T), true
		}
	}

	var zero T
	return zero, false
}

// Komut seçeneklerinden verilen anahtardaki değeri döndürür.
func (c *CommandLib) GetOptionBool(key string) (val bool, ok bool) {
	if key == "enabled" {
		return c.Options.Enabled, true
	}
	if key == "ephemeral" {
		return c.Options.Ephemeral, true
	}
	return getExtra[bool](c, key)
}

func (c *CommandLib) GetOptionInt(key string) (val int, ok bool) {
	if key == "cooldown" {
		return c.Options.Cooldown, true
	}
	return getExtra[int](c, key)
}

func (c *CommandLib) GetOptionString(key string) (val string, ok bool) {
	if key == "category" {
		return c.Options.Category, true
	}
	return getExtra[string](c, key)
}

func (c *CommandLib) GetAutoActions() []AutoAction {
	return c.Options.AutoActions
}

func (c *CommandLib) GetSlashArgs() []SlashArg {
	return c.Options.Slash
}

// Slash alanındaki Extra alanından verilen anahtardaki değeri döndürür.
func GetSlashExtraMap[T any](c *CommandLib, slashName, key string) (map[string]T, bool) {
	// İlgili SlashArg'ı bul
	var target *SlashArg
	for _, s := range c.Options.Slash {
		if s.Name == slashName {
			target = &s
			break
		}
	}

	if target == nil {
		return nil, false
	}

	// Extra içeriğini al
	v, ok := target.Extra[key]
	if !ok {
		return nil, false
	}

	// Tip dönüşümü yap
	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}

	result := make(map[string]T)
	for k, val := range m {
		if typedVal, ok := val.(T); ok {
			result[k] = typedVal
		}
	}

	return result, true
}
