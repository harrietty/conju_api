package verbs

// Infinitives is an array of string infinitive verbs
type Infinitives []string

// LanguageData contains an array of available pronouns and the Verbs struct
type LanguageData struct {
	Pronouns []string `json:"pronouns"`
	Verbs    Verbs    `json:"verbs"`
}

// Verbs is an array of Verb structs
type Verbs struct {
	Basic []Verb `json:"basic"`
}

// Verb represents a single verb with its translations conjugations
type Verb struct {
	Infinitive   string       `json:"infinitive"`
	Translations []string     `json:"translations"`
	Type         []string     `json:"type"`
	Conjugations Conjugations `json:"conjugations"`
}

// Conjugations is a struct of conjugations for a specific verb
type Conjugations struct {
	Present     []string `json:"present"`
	Preterite   []string `json:"preterite"`
	Imperfect   []string `json:"imperfect"`
	Conditional []string `json:"conditional"`
	Future      []string `json:"future"`
}

// EnglishData is a struct of English verb data
type EnglishData struct {
	Pronouns []string     `json:"pronouns"`
	Verbs    EnglishVerbs `json:"verbs"`
}

// EnglishVerbs is a struct containing different types of English verbs
type EnglishVerbs struct {
	Basic map[string]EnglishVerb `json:"basic"`
}

// EnglishVerb represents a single English verb and its conjugations
type EnglishVerb struct {
	Present           []string `json:"present"`
	Preterite         []string `json:"preterite"`
	PresentParticiple string   `json:"present_participle"`
	PastParticiple    string   `json:"past_participle"`
	Root              string   `json:"root"`
}
