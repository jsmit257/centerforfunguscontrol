#!/bin/bash

if which psql >/dev/null 2>&1; then
  psql "postgresql://postgres:root@${HUAUTLA_HOST:=NOTlocalhost}:${HUAUTLA_PORT:=5432}" huautla <<-EOF
    \c huautla
    update strains set generation_uuid = null;
    delete from notes;
    delete from photos;
    delete from sources;
    delete from generations;
    delete from events;
    delete from lifecycles;
    delete from strain_attributes;
    delete from strains;
    delete from substrate_ingredients;
    delete from substrates;
    delete from vendors where uuid != 'localhost';
EOF
fi

files=(
  ./tests/system/main_test.go
  ./tests/system/vendor_test.go
  ./tests/system/ingredient_test.go
  ./tests/system/substrate_test.go
  ./tests/system/substrateingredient_test.go
  ./tests/system/strain_test.go
  ./tests/system/strainattribute_test.go
  ./tests/system/lifecycle_test.go
  ./tests/system/generation_test.go
  ./tests/system/event_test.go
  ./tests/system/source_test.go
)

go get net/http

go test "${files[@]}"
