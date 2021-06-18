.PHONY: help build test deps generate clean

web_build:
	yarn install
	npm run build
	cp -r build/* static

web_clean:
	rm -rf build
	rm -rf node_modules
	rm -rf static/static
	find static ! -name '.gitignore' ! -name 'README.md'  -type f -exec rm -rf {} +