#!/bin/bash

# launch the bot
cd ../ && go build && ./pleasantbot &

# launch React front-end
cd src/gui && npm start
