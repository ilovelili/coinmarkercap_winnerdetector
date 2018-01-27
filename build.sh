#!/bin/sh
# makefile
cd src && make
# run ansible
cd ../ansible && ./playbooks/deployoncron.yml