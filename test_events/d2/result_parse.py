# Python code to
# demonstrate readlines()

import re

# writing to file
output = open("2019_d2_boys_events.txt", "w")

# Using readlines()
file1 = open("2019_d2_boys.txt", "r")
Lines = file1.readlines()

# Strips the newline character
for line in Lines:
    # PDF parser
    # line = line.strip()
    # # cols          0   1     2     3       4
    # # columns:  place, first, last, grade, school,...., time, score
    # cols = line.split(" ")
    # fields = [cols[0]]
    # fields.append(str(int(cols[0]) * 10))  # add bib
    # fields.append(cols[2])  # last
    # fields.append(cols[1])  # first
    # fields.append(cols[3])  # grade
    # nameCols = 1
    # if len(cols) > 7:
    #     nameCols = len(cols) - 7 + 1

    # # get the variable number of names
    # fields.append(" ".join(cols[4 : 4 + nameCols]))
    # # get the last 2 cols
    # fields.extend(cols[-2:])

    # old text parser
    line = line.strip().replace("#", "").replace(",", "")
    line = re.sub("[ ]+", " ", line)
    cols = line.split(" ")
    # cols          0    1     2     3      4       5
    # columns:  place, bib, last, first, grade, school,...., avg, time, score
    # school name is 5 - ?, 9 is min, so if len is 10, 2 part school name
    fields = cols[0:5]
    nameCols = 1
    if len(cols) > 9:
        nameCols = len(cols) - 9 + 1

    # get the variable number of names
    fields.append(" ".join(cols[5 : 5 + nameCols]))
    # get the last 3 cols
    fields.extend(cols[-2:])
    print("|".join(fields))
    output.write("|".join(fields) + "\n")
