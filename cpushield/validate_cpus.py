#!/usr/bin/env python3
import sys

total = int(sys.argv[1])
for p in sys.argv[2].split(","):
    p = p.strip()
    if not p: continue
    highest_requested = int(p.split("-")[1]) if "-" in p else int(p)
    if highest_requested >= total:
        print(highest_requested, end="")
        sys.exit(1)
