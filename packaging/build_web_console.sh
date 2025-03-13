#!/usr/bin/env bash

# set web console target dir
WEB_CONSOLE_TARGET_DIR="web"
# set web console git repo
WEB_CONSOLE_REPO="git@github.com:lindb/lin.git"
# build web console tmp dir
WEB_CONSOLE_TMP_DIR="web_console_build_temp"


function build_web_console() {
	echo "build web console start ..."
  git clone $WEB_CONSOLE_REPO $WEB_CONSOLE_TMP_DIR

	cd $WEB_CONSOLE_TMP_DIR
	#FIXME: checkout release tag/branch

	pnpm install
	pnpm build --filter @lindata/lindb

	cd ..
	find $WEB_CONSOLE_TARGET_DIR -mindepth 1 ! -name 'README.md' -type f -exec rm -rf {} +
	cp -r $WEB_CONSOLE_TMP_DIR/apps/lindb/dist/* $WEB_CONSOLE_TARGET_DIR

	echo "clean web console tmp dir..."
	rm -rf $WEB_CONSOLE_TMP_DIR
}

build_web_console 
