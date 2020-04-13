#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import urllib.request as r
import re

url = "https://ekstrabladet.dk/side9/"
regex = r"src=\"(/incoming/.*IMAGE_ALTERNATES/p)(\d+)([\w/-]+)\""

with r.urlopen(url) as html:
    g = (re.findall(regex, html.read().decode('utf-8'), re.MULTILINE))[0]
    print("https://ekstrabladet.dk{}1600{}".format(g[0], g[2]))
