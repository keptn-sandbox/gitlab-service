#!/bin/env bash

export UUID=$(uuidgen)
export TESTID=$RANDOM
envsubst < configuration-change.http > tmp_configuration-change.http
cat tmp_configuration-change.http
keptn send event -f tmp_configuration-change.http
rm tmp_configuration-change.http
