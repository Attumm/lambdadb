# -*- coding: utf-8 -*-
"""
Load first rows from csv, ask some questions
and generate a models.go to jumpstart
your lambda_db project for the given csv file

models.go contains all the field information
and functions of rows in your data.

- Repeated option to store repeated
  values in a map and each individual items
  only stores uint16 reference to map key.

- BitArray option which is like Repeated
  value but also creates a map[key]bitmap for all
  items containing field value. Makes it possible
  to do fast 'match' lookups.


python create_model.py your.csv
"""

import csv
import sys
import os

from re import sub
from jinja2 import Environment, FileSystemLoader

import yaml


if '-f' in sys.argv:
    filename = str(sys.argv[sys.argv.index('-f')+1])
else:
    filename = "items.csv"

if '-c' in sys.argv:
    config = str(sys.argv[sys.argv.index('-c')+1])
else:
    config = "config.yaml"

with open(filename) as f:
    reader = csv.DictReader(f)
    row = dict(next(reader))

cfg = {}

if os.path.isfile(config):
    with open(config, 'r') as stream:
        cfg = yaml.load(stream)['model']


env = Environment(
    loader=FileSystemLoader('./templates'),
)

# keep track of all column names and all original names in csv
allcolumns = []
allcolumns_org = []
repeated = []
repeated_org = []
bitarray = []
bitarray_org = []
unique = []
unique = []
unique_org = []
ignored = []
ignored_org = []
geocolumns = []
geocolumns_org = []


def gocamelCase(string):
    """convert string to camelCase

    woning_type -> WoningType
    """
    string = sub(r"(_|-)+", " ", string).title().replace(" ", "")
    return string


# ask some questions about columns.
index = 0
for k in row.keys():

    # go camelcase column names
    kc = gocamelCase(k)

    options = ['r', 'u', 'i', 'g', 'b']
    while True:

        action = None

        if cfg.get(k):
            print(f"reading from config {k} {cfg[k]}")
            action = cfg[k]
        else:
            # keep asking for valid input
            q1 = (
                "(R)epeated value? has less then (2^16=65536) option.",
                "(B)itarray, repeated column optimized for fast match.",
                "(U)nique, (G)eo lat/lon point or (I)gnore ? r/b/u/g/i?."
            )
            action = input(f"idx:{index} is {k} {q1}")  # noqa

        if action == '':
            print(f"pick one from {options}")
            continue
        if action not in options:
            continue
        break

    cfg[k] = action

    if action == 'r':
        repeated.append(kc)
        repeated_org.append(k)
    elif action == 'u':
        unique.append(kc)
        unique_org.append(k)
    elif action == 'i':
        ignored.append(kc)
        ignored_org.append(k)
    elif action == 'g':
        geocolumns.append(kc)
        geocolumns_org.append(k)
        unique.append(kc)
        unique_org.append(k)
    elif action == 'b':
        # same as repeated  but with some extra bitarray stuff
        repeated.append(kc)
        repeated_org.append(k)
        bitarray.append(kc)
        bitarray_org.append(k)
    else:
        print('invalid input')
        sys.exit(-1)

    allcolumns.append(kc)
    allcolumns_org.append(k)
    index += 1

# ask for a index column
while True:
    index = None
    # keep asking for valid input
    if cfg.get('index'):
        index = cfg['index']
    else:
        index = input(f"which column is idx? 0 - {len(allcolumns) - 1} ")

    cfg['index'] = index

    try:
        index = int(index)

        if allcolumns[index] in ignored:
            print('Selected an ignored column for index')
            raise ValueError

        if -1 < index < len(allcolumns):
            break

    except ValueError:
        continue

    print('try again..')

# save answers in config file
with open(config, 'w') as f:
    dict_file = {'model': cfg}
    yaml.dump(dict_file, f)
    print(f'saved answers in config {config}')


# setup initial data structs for each repeated column
initRepeatColumns = []
repeatColumnNames = []
loadRepeatColumnNames = []

for columnName in repeated:
    initRow = f"\t {columnName} = NewReapeatedColumn()\n"
    initRepeatColumns.append(initRow)

    repeatRow = f"\t {columnName} \n"
    repeatColumnNames.append(repeatRow)

    loadRow = f"\t {columnName} = m.{columnName} \n"
    loadRepeatColumnNames.append(loadRow)


# setup initial data structs for each bitarray column
initBitarrays = []
for columnName in bitarray:
    onerow = f"\t {columnName}Items = make(fieldItemsMap)\n"
    initBitarrays.append(onerow)

# create bitarrays with item labels for column values.
bitArrayStores = []
for c1, c2 in zip(bitarray, bitarray_org):
    onerow = f'\tSetBitArray("c2", i.{c1}, i.Label)\n'
    bitArrayStores.append(onerow)


# create ItemFull struct fields
columnsItemIn = []
jsonColumn = env.get_template('itemFullColumn.jinja2')
for c1, c2 in zip(allcolumns, allcolumns_org):
    onerow = jsonColumn.render(c1=c1, c2=c2)
    columnsItemIn.append(onerow)

# create ItemFull struct fields
columnsItemOut = []
jsonColumn = env.get_template('itemFullColumn.jinja2')
for c1, c2 in zip(allcolumns, allcolumns_org):

    if c1 in ignored:
        continue

    onerow = jsonColumn.render(c1=c1, c2=c2)
    columnsItemOut.append(onerow)

# create Item struct fields
columnsItem = []
for c1, c2 in zip(allcolumns, allcolumns_org):

    if c1 in ignored:
        continue

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
    shrinkVars.append(
        shrinkvartemplate.render(column=c, bitarray=c in bitarray))
    shrinkItems.append(f"\t {c}.Store(i.{c})\n")


# create the actual shrinked/expand Item fields.
shrinkItemFields = []
expandItemFields = []

for c in allcolumns:

    if c in ignored:
        continue

    if c in repeated:
        # string to unint
        shrinkItemFields.append(f"\t\t{c}.GetIndex(i.{c}),\n")
        # unint back to string
        expandItemFields.append(f"\t\t{c}.GetValue(i.{c}),\n")
    else:
        shrinkItemFields.append(f"\t\ti.{c},\n")
        expandItemFields.append(f"\t\ti.{c},\n")


# ItemIn Columns
inColumns = []
for c in allcolumns_org:
    inColumns.append(f'\t\t"{c}",\n')

# ItemOut Columns
outColumns = []
for cc, c in zip(allcolumns, allcolumns_org):
    # cc CamelCaseColumn.
    if cc in ignored:
        continue
    outColumns.append(f'\t\t"{c}",\n')

# create column filters.
# match, startswith, contains etc

columnFilters = []
filtertemplate = env.get_template("filters.jinja2")

for c in allcolumns:
    if c in ignored:
        continue

    lookup = f"i.{c}"
    if c in repeated:
        lookup = f"{c}[i.{c}]"

    txt = filtertemplate.render(column=c, lookup=lookup)
    columnFilters.append(txt)

registerFilters = []
rtempl = env.get_template('registerFilters.jinja2')
# register filters
for c, co in zip(allcolumns, allcolumns_org):
    if c in ignored:
        continue
    txt = rtempl.render(co=co, columnName=c, bitarray=c in bitarray)
    registerFilters.append(txt)

sortColumns = []
sortTemplate = env.get_template('sortfunc.jinja2')

# create sort functions
for c, co in zip(allcolumns, allcolumns_org):
    if c in ignored:
        continue

    c1 = f"items[i].{c} < items[j].{c}"
    c2 = f"items[i].{c} > items[j].{c}"

    if c in repeated:
        c1 = f"{c}[items[i].{c}] < {c}[items[j].{c}]"
        c2 = f"{c}[items[i].{c}] > {c}[items[j].{c}]"

    txt = sortTemplate.render(co=co, c1=c1, c2=c2)
    sortColumns.append(txt)


csv_columns = []
for c in allcolumns:
    csv_columns.append(f'\t"{c}",\n')


# Finally render the model.go template
modeltemplate = env.get_template('model.template.jinja2')
mapstemplate = env.get_template('modelmap.template.jinja2')

geometryGetter = '""'
print('GEOCOLUMNS: ' + " ".join(geocolumns))
if len(geocolumns) == 1:
    geometryGetter = f"Getters{geocolumns[0]}(&i)"

output = modeltemplate.render(
    #initRepeatColumns=''.join(initRepeatColumns),
    columnsItemIn=''.join(columnsItemIn),
    columnsItemOut=''.join(columnsItemOut),
    columnsItem=''.join(columnsItem),
    # shrinkVars=''.join(shrinkVars),
    shrinkItems=''.join(shrinkItems),
    shrinkItemFields=''.join(shrinkItemFields),
    expandItemFields=''.join(expandItemFields),
    csv_columns=''.join(csv_columns),
    inColumns=''.join(inColumns),
    outColumns=''.join(outColumns),
    columnFilters=''.join(columnFilters),
    registerFilters=''.join(registerFilters),
    sortColumns=''.join(sortColumns),
    indexcolumn=allcolumns[index],
    geometryGetter=geometryGetter,
    bitArrayStores=''.join(bitArrayStores),
)

f = open('model.go', 'w')
f.write(output)
f.close()
print('saved in model.go')
print('!!NOTE!! edit the default search filter')


mapsoutput = mapstemplate.render(
    initRepeatColumns=''.join(initRepeatColumns),
    repeatColumnNames = ''.join(repeatColumnNames),
    loadRepeatColumnNames = ''.join(loadRepeatColumnNames),
    initBitarrays=''.join(initBitarrays),
    shrinkVars=''.join(shrinkVars),

)

f = open('modelmaps.go', 'w')
f.write(mapsoutput)
f.close()
print('model hashmaps  saved in modelmaps.go')


