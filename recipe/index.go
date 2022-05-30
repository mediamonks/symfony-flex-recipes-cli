package recipe

type Index struct {
	Aliases   map[string]string   `json:"aliases"`
	Recipes   map[string][]string `json:"recipes"`
	Branch    string              `json:"branch"`
	IsContrib bool                `json:"is_contrib"`
	Links     Links               `json:"_links"`
}

type Links struct {
	Repository                      string `json:"repository"`
	OriginTemplate                  string `json:"origin_template"`
	RecipeTemplate                  string `json:"recipe_template"`
	RecipeTemplateRelative          string `json:"recipe_template_relative"`
	ArchivedRecipesTemplate         string `json:"archived_recipes_template"`
	ArchivedRecipesTemplateRelative string `json:"archived_recipes_template_relative"`
}
