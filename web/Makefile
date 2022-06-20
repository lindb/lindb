.PHONY: help build test deps generate clean

web_build: web_clean
	yarn install
	yarn build
	mkdir -p static && cp -r build/* static

web_clean:
	rm -rf build
	rm -rf static/assets
	find static ! -name '.gitignore' ! -name 'README.md'  -type f -exec rm -rf {} +

web_clean_all:
	rm -rf build
	rm -rf node_modules
	find static ! -name '.gitignore' ! -name 'README.md'  -type f -exec rm -rf {} +
