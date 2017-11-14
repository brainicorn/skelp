#!/usr/bin/env bash

first_arg=$1
shift
echo "Greetings $@" > $first_arg
