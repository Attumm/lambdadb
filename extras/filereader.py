import csv


def supported_fileformats():
    return list(FILE_FORMATS.keys())


def create_reader(f, file_format):
    reader_data = FILE_FORMATS[file_format]
    return reader_data["func"](f, **reader_data.get("args", {}))


FILE_FORMATS = {
    "csv": {"func": csv.DictReader},
    "tsv": {"func": csv.DictReader, "args": {"delimiter": "\t", "quotechar": '"'}},
}

