#!/bin/sh

(for workers in 10 30 100 ; do for rate in 250 275 290 300 310 325 350 ; do go run . $workers $rate ; done ; done) > plot.txt
python3 plot.py
