sleep 1
python3 extras/ingestion.py  -f examples/movies_subset.tsv -format tsv -dbhost 127.0.0.1:8128

curl 127.0.0.1:8128/mgmt/save/
curl 127.0.0.1:8128/mgmt/rm/
curl 127.0.0.1:8128/mgmt/load/
