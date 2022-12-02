#!/usr/bin/env python3

import sys

elfPacks = []
currPack = 0

for line in sys.stdin:
    line = line.strip()
    if line == "":
        elfPacks.append(currPack)
        currPack = 0
    else:
        currPack += int(line)

if currPack != 0:
    elfPacks.append(currPack)

elfPacks.sort()

print("The largest pack is: {}".format(elfPacks[len(elfPacks) - 1]))

topCount = 3
total = 0

for _, val in enumerate(elfPacks[len(elfPacks) - topCount:]):
    total += val

print("The top {#010d} elves have a total of {}".format(topCount, total))
