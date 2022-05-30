package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	json2 "recipes-cli/json"
	"recipes-cli/recipe"
	"strings"
)

var generateCommand = &cobra.Command{
	Use:     "generate",
	Short:   "Generates the recipes index and manifests",
	Example: "recipes generate",
	RunE: func(cmd *cobra.Command, args []string) error {
		index := &recipe.Index{
			Aliases:   map[string]string{},
			Recipes:   make(map[string][]string),
			Branch:    "tree",
			IsContrib: true,
			Links: recipe.Links{
				Repository:                      fmt.Sprintf("github.com/%s", os.Getenv("GITHUB_REPOSITORY")),
				OriginTemplate:                  fmt.Sprintf("{package}:{version}@github.com/%s:tree", os.Getenv("GITHUB_REPOSITORY")),
				RecipeTemplate:                  fmt.Sprintf("https://raw.githubusercontent.com/%s/tree/{package_dotted}.{version}.json", os.Getenv("GITHUB_REPOSITORY")),
				RecipeTemplateRelative:          "{package_dotted}.{version}.json",
				ArchivedRecipesTemplate:         fmt.Sprintf("https://raw.githubusercontent.com/%s/main/archived/{package_dotted}/{ref}.json", os.Getenv("GITHUB_REPOSITORY")),
				ArchivedRecipesTemplateRelative: "archived/{package_dotted}/{ref}.json",
			},
		}

		// git ls-tree HEAD */*/* | ./recipes generate
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			tokens := strings.Split(strings.TrimSpace(scanner.Text()), "\t")
			tree := strings.Split(tokens[0], " ")
			ref := tree[2]
			pkg := tokens[1:2][0]

			manifestBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/manifest.json", pkg))
			if err != nil {
				return fmt.Errorf("unable to read %s/manifest.json: %w", pkg, err)
			}

			manifest := &recipe.Manifest{}
			if err := json.Unmarshal(manifestBytes, manifest); err != nil {
				return err
			}

			version := pkg[1+strings.LastIndex(pkg, "/"):]
			pkg = pkg[:len(pkg)-len(version)-1]

			for _, alias := range manifest.Aliases {
				index.Aliases[alias] = pkg
			}

			// We don't want the aliases in the recipe manifest
			manifest.Aliases = nil

			if err := generatePackageJson(pkg, version, ref, manifest); err != nil {
				return err
			}

			// Add recipe's
			if _, ok := index.Recipes[pkg]; !ok {
				index.Recipes[pkg] = make([]string, 0)
			}

			index.Recipes[pkg] = append(index.Recipes[pkg], version)
		}

		indexJson, err := json2.Marshal(index)
		if err != nil {
			return fmt.Errorf("unable to marshal index.json: %w", err)
		}

		if err := ioutil.WriteFile("index.json", indexJson, os.ModePerm); err != nil {
			return fmt.Errorf("unable to save index.json: %w", err)
		}

		return nil
	},
}

func generatePackageJson(pkg string, version string, ref string, manifest *recipe.Manifest) error {
	files := make(map[string]*recipe.File, 0)
	err := filepath.Walk(fmt.Sprintf("%s/%s", pkg, version), func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || info.Name() == "manifest.json" {
			return nil
		}

		if info.Name() == "post-install.txt" {
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("unable to read post-install.txt: %w", err)
			}

			manifest.PostInstallOutput = strings.Split(strings.TrimSpace(strings.ReplaceAll(string(contents), "\r", "")), "\n")
			return nil
		}

		if info.Name() == "Makefile" {
			// @todo also parse Makefile
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file, %s: %w", path, err)
		}

		fileContents := strings.Split(string(contents), "\n")

		stats, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("unable to retrieve file stats for, %s: %w", path, err)
		}

		files[strings.Replace(path, fmt.Sprintf("%s/%s/", pkg, version), "", 1)] = &recipe.File{
			Contents:   fileContents,
			Executable: IsExecutable(stats.Mode()),
		}

		return nil
	})

	if err != nil {
		return err
	}

	rcpe := &recipe.Recipe{
		Manifests: map[string]*recipe.ManifestEntry{},
	}

	if _, ok := rcpe.Manifests[pkg]; !ok {
		rcpe.Manifests[pkg] = &recipe.ManifestEntry{
			Manifest: manifest,
			Files:    files,
			Ref:      ref,
		}
	}

	// Also create an archived recipe
	packageNameDotted := strings.ReplaceAll(pkg, "/", ".")
	archivePath := fmt.Sprintf("./archived/%s", packageNameDotted)

	recipeBytes, err := json2.Marshal(rcpe)
	if err != nil {
		return fmt.Errorf("unable to marshal recipe: %w", err)
	}

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		if err := os.MkdirAll(archivePath, os.ModePerm); err != nil {
			return fmt.Errorf("unable to make archive dir, %s: %w", archivePath, err)
		}
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.json", archivePath, ref), recipeBytes, os.ModePerm); err != nil {
		return fmt.Errorf("unable to write archive file: %s/%s.json: %w", archivePath, ref, err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s.%s.json", packageNameDotted, version), recipeBytes, os.ModePerm); err != nil {
		return fmt.Errorf("unable to write recipe file: %.%s.json: %w", packageNameDotted, version, err)
	}

	return nil
}

func IsExecutable(mode os.FileMode) bool {
	return mode&0111 == 0111
}
