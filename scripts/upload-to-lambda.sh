#!/bin/bash -ex

# todo cd to the right folder

./build.sh
zip gobot.zip gobot

aws lambda update-function-code \
--function-name gobot \
--zip-file fileb://gobot.zip
rm -f gobot.zip