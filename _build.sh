#!/usr/bin/env bash

goxc -d=./bin -bc="linux,386,amd64 darwin,386,amd64"
rm -rf debian