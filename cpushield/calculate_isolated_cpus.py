#!/usr/bin/env python3
import sys

total = int(sys.argv[1])
cpu_str = sys.argv[2]
cpus = set()

for p in cpu_str.split(","):
    p = p.strip()
    if not p: continue
    if "-" in p:
        s, e = map(int, p.split("-"))
        cpus.update(range(s, e+1))
    else:
        cpus.add(int(p))

iso = sorted(list(set(range(total)) - cpus))
res = []

if iso:
    s = e = iso[0]
    for n in iso[1:]:
        if n == e + 1:
            e = n
        else:
            res.append(str(s) if s == e else f"{s}-{e}")
            s = e = n
    res.append(str(s) if s == e else f"{s}-{e}")

print(",".join(res), end="")
