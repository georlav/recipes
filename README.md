![Tests](https://github.com/georlav/migrate/workflows/Tests/badge.svg?branch=master)
![GolangCI](https://github.com/georlav/migrate/workflows/GolangCI/badge.svg?branch=master)

# Recipes
A simple app that fetches puppy recipes from various providers and prints them to the standard output. Use those 
recipes to cook for your puppy (not your puppy).

## Current supported providers
http://www.recipepuppy.com/api

## Install only the command line tool
```bash
go get github.com/georlav/recipes/cmd/recipes/...
```

## Makefile
To run the app
```bash
make run
```

To run the tests
```bash
make test
```

To produce a build
```bash
make build
```

## Configuration
Most options can be changed from config.json

## Authors
* **georlav** - *Initial work* - [georlav](https://github.com/georlav)

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details



