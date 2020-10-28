# -*- coding: utf-8 -*-
"""
Load first rows from csv, ask some questions
and generate a models.go to jumpstart
your project for the given csv file.

Much morge memory efficient then v1 because repeated
values are now stored in a map and each individual item
only stores uint16 reference.

python create_model.py your.csv
"""

import csv
import sys

from re import sub
from jinja2 import Environment, FileSystemLoader

if '-f' in sys.argv:
    filename = str(sys.argv[sys.argv.index('-f')+1])
else:
    filename = "items.csv"

with open(filename) as f:
    reader = csv.DictReader(f)
    row = dict(next(reader))

env = Environment(
    loader=FileSystemLoader('./templates'),
)

# keep track of all column names org are original names
allcolumns = []
repeated = []
repeated_org = []
unique = []
unique_org = []


def gocamelCase(string):
    """convert string to camelCase

    woning_type -> WoningType
    """
    string = sub(r"(_|-)+", " ", string).title().replace(" ", "")
    return string


# ask some questions about columns.
index = 0
for k in row.keys():
    kc = gocamelCase(k)

    while True:
        # keep asking for valid input
        q1 = "a repeated value? has less then (2^16=65536) values? Y/n?"
        yesno = input(f"idx:{index} is {k} {q1}")  # noqa
        if yesno == '':
            yesno = 'y'
        if yesno not in ['y', 'n']:
            continue
        break

    if yesno == 'y':
        repeated.append(kc)
        repeated_org.append(k)
    else:
        unique.append(kc)
        unique_org.append(k)

    allcolumns.append(kc)
    index += 1

# ask for a index column
while True:
    # keep asking for valid input
    index = input(f"which column is idx? 0 - {len(allcolumns) - 1} ")
    try:
        index = int(index)
        if index < len(allcolumns):
            break
    except ValueError:
        continue
    print('try again..')

# setup initial data structs for each repeated column
initRepeatColumns = []
initColumntemplate = env.get_template('initColumn.template.jinja2')

for c in repeated:
    initRepeatColumns.append(initColumntemplate.render(columnname=c))

# create ItemFull struct fields
columnsItemFull = []
jsonColumn = env.get_template('itemFullColumn.jinja2')
for c1, c2 in zip(allcolumns, row.keys()):
    onerow = jsonColumn.render(c1=c1, c2=c2)
    columnsItemFull.append(onerow)

# create Item struct fields
columnsItem = []
for c1, c2 in zip(allcolumns, row.keys()):
    onerow = f"\t{c1}  string\n"
    if c1 in repeated:
        onerow = f"\t{c1}    uint16\n"
    columnsItem.append(onerow)


# create Shrink code for repeated fields
# where we map uint16 to a string value.
shrinkVars = []
shrinkItems = []
shrinkvartemplate = env.get_template('shrinkVars.jinja2')
shrinktemplate = env.get_template('shrinkColumn.jinja2')
for c in repeated:
    shrinkVars.append(shrinkvartemplate.render(column=c))
    shrinkItems.append(shrinktemplate.render(column=c))


# create the actual shrinked/expand Item fields.
shrinkItemFields = []
expandItemFields = []

for c in allcolumns:
    if c in repeated:
        # string to unint
        shrinkItemFields.append(f"\t\t{c}IdxMap[i.{c}],\n")
        # unint back to string
        expandItemFields.append(f"\t\t{c}[i.{c}],\n")
    else:
        shrinkItemFields.append(f"\t\ti.{c},\n")
        expandItemFields.append(f"\t\ti.{c},\n")


originalColumns = []
for c in row.keys():
    originalColumns.append(f'\t\t"{c}",\n')

# create column filters.
# match, startswith, contains etc

columnFilters = []
filtertemplate = env.get_template("filters.jinja2")

for c in allcolumns:
    lookup = f"i.{c}"
    if c in repeated:
        lookup = f"{c}[i.{c}]"

    txt = filtertemplate.render(column=c, lookup=lookup)
    columnFilters.append(txt)

registerFilters = []
rtempl = env.get_template('registerFilters.jinja2')
# register filters
for c, co in zip(allcolumns, row.keys()):
    txt = rtempl.render(co=co, column=c)
    registerFilters.append(txt)


sortColumns = []
sortTemplate = env.get_template('sortfunc.jinja2')

# create sort functions
for co, c in zip(row.keys(), allcolumns):

    c1 = f"items[i].{c} < items[j].{c}"
    c2 = f"items[i].{c} > items[j].{c}"

    if c in repeated:
        c1 = f"{c}[items[i].{c}] < {c}[items[j].{c}]"
        c2 = f"{c}[items[i].{c}] > {c}[items[j].{c}]"

    txt = sortTemplate.render(co=co, c1=c1, c2=c2)
    sortColumns.append(txt)


csv_columns = []
for c in row.keys():
    csv_columns.append(f'\t"{c}",\n')


# Finally render the model.go template
modeltemplate = env.get_template('model.template.jinja2')

output = modeltemplate.render(
    initRepeatColumns=''.join(initRepeatColumns),
    columnsItemFull=''.join(columnsItemFull),
    columnsItem=''.join(columnsItem),
    shrinkVars=''.join(shrinkVars),
    shrinkItems=''.join(shrinkItems),
    shrinkItemFields=''.join(shrinkItemFields),
    expandItemFields=''.join(expandItemFields),
    csv_columns=''.join(csv_columns),
    originalColumns=''.join(originalColumns),
    columnFilters=''.join(columnFilters),
    registerFilters=''.join(registerFilters),
    sortColumns=''.join(sortColumns),
    indexcolumn=allcolumns[index]
)

f = open('model.go', 'w')
f.write(output)
f.close()

print('saved in model.go')
print('!!NOTE!! edit the default search filter')
