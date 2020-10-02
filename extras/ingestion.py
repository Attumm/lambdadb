
import sys
import json
import requests
import csv

csv.field_size_limit(sys.maxsize)

if __name__ == "__main__":
    produce = bool(sys.argv[sys.argv.index('-p')+1]) if '-p' in sys.argv else False 
    start_worker = int(sys.argv[sys.argv.index('-w')+1]) if '-w' in sys.argv else 0
    produce_http = bool(sys.argv[sys.argv.index('-phttp')+1]) if '-phttp' in sys.argv else False
    filename = str(sys.argv[sys.argv.index('-f')+1]) if '-f' in sys.argv else "items.csv"
    buffer_size = int(sys.argv[sys.argv.index('-b')+1]) if '-b' in sys.argv else 100000
    http_db_host = str(sys.argv[sys.argv.index('-dbhost')+1]) if '-dbhost' in sys.argv else "http://127.0.0.1:8000/add/"


    lines = []
    with open(filename) as f:
        reader = csv.DictReader(f)
        for i, row in enumerate(reader):
            parsed = {k.lower(): str(v) for k, v in row.items()}
            lines.append(parsed)
            if i % buffer_size == 0:
                requests.post(http_db_host, json=lines)
                lines = []

        requests.post(http_db_host, json=lines)

