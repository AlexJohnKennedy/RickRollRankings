#!/bin/bash

# Figure out the directory this script is sitting in, so we can set the working directory here.
# This way we will be able to find the .env and .py files we need, regardless of where this script is
# being run from.
scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")
cd $scriptDir

# Do some stupid shit where you parse variables inside an .env file with cat, and dump them into an 'env'
# command which runs the .py file, providing the environment vars I need. This is all very hacky and terrible
# but I'm just playing around atm anyway.
env $(cat ./redditApiKeys.env | tr "\\n" " ") python3 ./RickRollCrawler.py