#!/bin/sh

if [ "$#" -ne 2 ]; then
  echo "usage: init.sh YEAR DAY"
  exit 1
fi

cat template.go.tmpl | sed "s/__YEAR__/$1/g; s/__DAY__/$2/g" > day$2/a/main.go
cat template.go.tmpl | sed "s/__YEAR__/$1/g; s/__DAY__/$2/g" > day$2/b/main.go
