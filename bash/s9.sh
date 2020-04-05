#!/bin/sh
curl -qks -XGET "https://ekstrabladet.dk/side9/"|sed -n 's#.*src="\(.*NATES/p*\)[0-9]*\(/\S*\)".*#https://ekstrabladet.dk\11600\2#gp'
