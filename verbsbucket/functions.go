package verbsbucket

// VerbsFileGetter gets a verbs file from the verbs bucket
type VerbsFileGetter func(string) ([]byte, error)
