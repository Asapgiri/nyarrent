#!/bin/bash

transmission-remote --list | sed -e '1d;$d;s/^ *//' | cut -f1 -d' ' | sed s/^/transmission-remote\ -t\ / | sed s/$/\ -rad/

