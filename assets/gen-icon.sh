#!/usr/bin/env bash

magick -size 80x80 xc:none -fill white \
    -draw "circle 40,50 40,30" \
    -draw "polygon 40,10 22,45 58,45" \
    icon.png