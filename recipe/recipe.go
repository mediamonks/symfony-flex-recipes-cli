package recipe

type Recipe struct {
	Manifests map[string]*ManifestEntry `json:"manifests"`
}

type ManifestEntry struct {
	Manifest *Manifest        `json:"manifest"`
	Files    map[string]*File `json:"files"`
	Ref      string           `json:"ref"`
}

type Manifest struct {
	Aliases           []string          `json:"aliases,omitempty"`
	ComposerScripts   map[string]string `json:"composer-scripts"`
	CopyFromRecipe    map[string]string `json:"copy-from-recipe"`
	PostInstallOutput []string          `json:"post-install-output"`
	Env               map[string]string `json:"env"`
}

type File struct {
	Contents   []string `json:"contents"`
	Executable bool     `json:"executable"`
}
