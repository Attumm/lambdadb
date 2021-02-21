### First version is going to assume everything is a string
### also known as string theory:p

### column with the name "value" or "index" will be used as index
### else the first column will be set as index, when index is enabled.
### this can be changed later in the generated model.go file

import csv
import sys
from filereader import create_reader, supported_fileformats


def create_struct(item):
    start = "type Item struct {\n"
    # TODO add type
    lines= [f'{k.capitalize()} string `json:"{k.lower()}"`' for k, v in item.items()]
    stop = "\n}\n"
    return start + "\n".join(lines) + stop 


def create_columns(item):
    start = """
    func (i Item) Columns() []string {
	return []string{
    """
    lines = [f'"{k.lower()}",' for k in item.keys()]
    stop = """\n}\n}"""
    return start + "\n".join(lines) + stop


def create_row(item):
    start = """
    func (i Item) Row() []string {
	return []string{
    """
    lines = [f"i.{k.capitalize()}," for k in item.keys()]
    stop = """\n}\n}"""
    return start + "\n".join(lines) + stop


def get_index_column(item):
    special_columns = ["value", "index"]
    for column in special_columns:
        if column in item:
            return column

    # we tried, let's return the first column
    n = iter(item.keys())
    return next(n)


def create_getindex(item):
    index_column = get_index_column(item)
    start = """
    func (i Item) GetIndex() string {
	return """
    middle = f"i.{index_column.capitalize()}"
    stop = """\n}"""
    return start + middle + stop


def create_filter_contains(column):
    return (
            f"func Filter{column.capitalize()}Contains(i *Item, s string) bool"  + "{" + "\n"
            f"return strings.Contains(i.{column.capitalize()}, s)"
            "\n" + "}"
    )

def create_filter_startswith(column):
    return (
            f"func Filter{column.capitalize()}StartsWith(i *Item, s string) bool"  + "{" + "\n"
            f"return strings.HasPrefix(i.{column.capitalize()}, s)"
            "\n" + "}"
    )

def create_filter_match(column):
    return (
            f"func Filter{column.capitalize()}Match(i *Item, s string) bool"  + "{" + "\n"
            f"return i.{column.capitalize()} ==  s"
            "\n" + "}"
    )


def create_getter(column):
    return (
            f"func Getters{column.capitalize()}(i *Item) string"  + "{" + "\n"
            f"return i.{column.capitalize()}"
            "\n" + "}"
    )


def create_reduce(column):
    return """
    func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}
    """

def create_init_register():
    return """
    RegisterFuncMap = make(registerFuncType)
    RegisterGroupBy = make(registerGroupByFunc)
    RegisterGetters = make(registerGettersMap)
    RegisterReduce = make(registerReduce)

    """

def create_register_match_func(column):
    return f'RegisterFuncMap["match-{column.lower()}"] = Filter{column.capitalize()}Match'


def create_register_contains_func(column):
    return f'RegisterFuncMap["contains-{column.lower()}"] = Filter{column.capitalize()}Contains'


def create_register_startswith_func(column):
    return f'RegisterFuncMap["startswith-{column.lower()}"] = Filter{column.capitalize()}StartsWith'


def create_register_getter(column):
    return f'RegisterGetters["{column.lower()}"] = Getters{column.capitalize()}'


def create_register_groupby(column):
    return f'RegisterGroupBy["{column.lower()}"] = Getters{column.capitalize()}'


def create_register_reduce(column):
    return 'RegisterReduce["count"] = reduceCount'


def create_grouped():
    return """
type GroupedOperations struct {
	Funcs   registerFuncType
	GroupBy registerGroupByFunc
	Getters registerGettersMap
	Reduce  registerReduce
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var RegisterGetters registerGettersMap
var RegisterReduce registerReduce
"""

def create_sortby_line_plus(column):
	return f'"{column.lower()}"' + ": func(i, j int) bool { return " + f"items[i].{column.capitalize()} < items[j].{column.capitalize()} " + " },"

def create_sortby_line_minus(column):
	return f'"-{column.lower()}"' + ": func(i, j int) bool { return " + f"items[i].{column.capitalize()} > items[j].{column.capitalize()} " + " },"

def create_sortby(row):
    start = """func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(int, int) bool{"""
    lines = []
    for k in row.keys():
        lines.append(create_sortby_line_plus(k))
        lines.append(create_sortby_line_minus(k))
        lines.append("\n")
    lines.append("}")
    end = """
    for _, sortFuncName := range sortingL {
    	sortFunc := sortFuncs[sortFuncName]
        sort.Slice(items, sortFunc)
                                }
        // TODO must be nicer way
        keys := []string{}
        for key := range sortFuncs {
              keys = append(keys, key)
        }

        return items, keys
        }"""
    return start + "\n".join(lines) + end

if __name__ == "__main__":

    filename = str(sys.argv[sys.argv.index('-f')+1]) if '-f' in sys.argv else "items.csv"
    file_format = str(sys.argv[sys.argv.index('-format')+1]) if '-format' in sys.argv else "csv"

    if file_format not in supported_fileformats():
        print(f"{file_format} not part of supported file formats {','.join(supported_fileformats())}")
        sys.exit()

    with open(filename) as f:
        reader = create_reader(f, file_format)
        row = dict(next(reader))

    print("package main")
    print()

    print("import (")
    print('"sort"')
    print('"strconv"')
    print('"strings"')
    print(")")
    print(create_struct(row))
    print()
    print(create_columns(row))
    print()
    print(create_row(row))
    print()
    print(create_getindex(row))
    print()

    print("// contain filters")
    for k in row.keys():
        print(create_filter_contains(k))

    print()
    print("// startswith filters")
    for k in row.keys():
        print(create_filter_startswith(k))

    print()
    print("// match filters")
    for k in row.keys():
        print(create_filter_match(k))

    print()
    print("// reduce functions")
    print(create_reduce(None))

    print()
    print("// getters")
    for k in row.keys():
        print(create_getter(k))
    print()


    print(create_grouped())
    print("func init() {")
    print(create_init_register())

    print()
    print("// register match filters")
    for k in row.keys():
        print(create_register_match_func(k))

    print()
    print("// register contains filters")
    for k in row.keys():
        print(create_register_contains_func(k))

    print()
    print("// register startswith filters")
    for k in row.keys():
        print(create_register_startswith_func(k))
    print()

    print()
    print("// register getters ")
    for k in row.keys():
        print(create_register_getter(k))
    print()

    print()
    print("// register groupby ")
    for k in row.keys():
        print(create_register_groupby(k))
    print()


    print()
    print("// register reduce functions")
    print(create_register_reduce(None))

    print("}")

    print(create_sortby(row))
    print()

