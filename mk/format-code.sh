#!/bin/bash
go fmt ./...
find -name '*.h' -o -name '*.c' \
  | grep -vE 'pcg_basic|siphash-20121104|uthash|zf_log' \
  | xargs clang-format-6.0 -i -style='{BasedOnStyle: Mozilla, ReflowComments: false}'
find . -path ./node_modules -prune -o \( -name '*.yaml' \) -print | xargs yamllint
node_modules/.bin/tslint --fix -p .