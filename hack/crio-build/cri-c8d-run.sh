#!/bin/bash

trap handle_exit EXIT

handle_exit() {
	pkill crio
}

crio --profile &

wait
