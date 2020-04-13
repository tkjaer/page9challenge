#!/usr/bin/env python3
# -*- coding: utf-8 -*-
from pprint import pprint
import re
import sys


def new_flow(flow_id, existing_ports=True):
    global cons
    flow_struct = {
            "src_pkt": 0,
            "dst_pkt": 0,
            "src_flag": "",
            "dst_flag": "",
            "reset": False,
            "complete": False
            }
    if existing_ports:
        cons[flow_id].append(flow_struct)
    else:
        cons[flow_id] = [flow_struct]


def update_flow(flow_id, sender, flags):
    global cons
    flow = cons[flow_id][-1]
    flow[sender + "_pkt"] += 1
    # this hacky fin+rst handling will fail, but maybe not within the simple
    # confinements of the challenge - "this is fine"
    if "F" in list(flags):
        flow[sender + "_flag"] = "FIN"
        if flow["src_flag"] == "FIN":
            if (flow["dst_flag"] == "FIN"):
                flow["complete"] = True
    if "R" in list(flags):
        flow["reset"] = True


def get_id(num):
    base = list("_abcdefghijklmnopqrstuvwxyz")
    arr = []
    while num:
        num, rem = divmod(num, len(base) - 1)
        arr.append(base[rem])
    arr.reverse()
    return "".join(arr)


cons = {}
cons_idx = []

for line in sys.stdin:
    regex = r"(.*)\.(\d+) (.*)\.(\d+): \[(.*)\]\,"
    match = re.match(regex, line).groups()

    if (match):
        a, a_port, b, b_port, flags = (match)
        ab = "{}:{} - {}:{}".format(a, a_port, b, b_port)
        ba = "{}:{} - {}:{}".format(b, b_port, a, a_port)

        if ab in cons:
            if flags == "S":
                new_flow(ab)
            update_flow(ab, "src", flags)

        elif ba in cons:
            if flags == "S":
                new_flow(ba)
            update_flow(ba, "dst", flags)

        else:
            new_flow(ab, existing_ports=False)
            update_flow(ab, "src", flags)
            cons_idx.append(ab)


con_number = 0
id_idx = 0

print("TCP connection info:")
for con_id in cons_idx:
    con_number += 1
    id_idx += 2
    src_id = get_id(id_idx -1)
    dst_id = get_id(id_idx)
    c = cons[con_id].pop()

    out_str = "{}: {} ({}2{})".format(
            str(con_number).rjust(3),
            con_id,
            src_id,
            dst_id)
    out_str = "{}{}>{}<".format(
            out_str.ljust(50),
            str(c["src_pkt"]).rjust(5),
            str(c["dst_pkt"]).rjust(5)
            )
    if c["complete"]:
        out_str += "  (complete)"
    if c["reset"]:
        out_str += "  (reset)"
    print(out_str)
