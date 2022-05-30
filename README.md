# Symfony Flex Recipes CLI

This cli tool helps generating private recipe repositories.

Branches:
- tree
- main

### Tree
mediamonks/{package}/{version}/  
-- config  
-- src  
-- manifest.json  
-- post-install.txt  

### Main
Contains the index.json and the recipe json files.

### Environment variables
- GITHUB_REPOSITORY

### Usage
While on the tree branch:  
```git ls-tree HEAD */*/* | ./recipes generate```  
Then:  
```
git switch master
git add *.json
git add archived/
git commit -m "updated recipes"
git push origin -f main
```