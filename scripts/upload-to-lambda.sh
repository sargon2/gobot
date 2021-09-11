#!/bin/bash -ex

cd "$(dirname "${BASH_SOURCE[0]}")"
cd ..

./build.sh
zip gobot.zip gobot

aws lambda update-function-code --function-name gobot --zip-file fileb://gobot.zip

rm -f gobot.zip
