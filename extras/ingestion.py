
import sys
import json
import requests
import csv

from filereader import create_reader, supported_fileformats

csv.field_size_limit(sys.maxsize)


def parse_command_line(flag, default_value, cast_to=str):
    flag_used = "-" + flag
    return cast_to(sys.argv[sys.argv.index(flag_used)+1]) if flag_used in sys.argv else default_value


if __name__ == "__main__":
    produce = parse_command_line(flag='p', default_value=False, cast_to=bool)
    start_worker = parse_command_line(flag='w', default_value=1, cast_to=int)
    produce_http = parse_command_line(flag='phttp', default_value=False, cast_to=bool)
    buffer_size = parse_command_line(flag='b', default_value=100000, cast_to=int)

    filename = parse_command_line(flag='f', default_value='items.csv', cast_to=str)
    file_format = parse_command_line(flag='format', default_value="csv", cast_to=str)
    http_db_host = parse_command_line(flag='dbhost', default_value="127.0.0.1:8128", cast_to=str)
    http_db_scheme = parse_command_line(flag='dbscheme', default_value="http://", cast_to=str)
    debug_mode = parse_command_line(flag='debug', default_value=False, cast_to=bool)
    link_to_add = "/mgmt/add/"

    #Just moby things, no walrus operation here.
    if file_format not in supported_fileformats():
        print(f"{file_format} not part of supported file formats {','.join(supported_fileformats())}")
        sys.exit()

    http_db_url = http_db_scheme + http_db_host + link_to_add
    lines = []
    with open(filename) as f:
        reader = create_reader(f, file_format)
        for i, row in enumerate(reader, start=1):
            parsed = {str(k).lower(): str(v).replace("\\", "") for k, v in row.items()}
            if debug_mode:
                print(i, parsed)
            lines.append(parsed)
            if i % buffer_size == 0:
                r = requests.post(http_db_url, json=lines)
                if r.status_code == 406:
                    print("LambdaDB is unable to process batch, try running this script with -debug flag")
                    exit(1)

                if debug_mode:
                    print(i % buffer_size, "status", r)
                lines = []

        r = requests.post(http_db_url, json=lines)
        if debug_mode:
            print("status", r)
        if r.status_code == 406:
            print("LambdaDB is unable to process batch, try running this script with -debug flag")
            exit(1)

