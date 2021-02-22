# LambdaDB
In memory database that uses filters to get the data you need.

Can be used for your needs by changing the models.go file to your needs.
Creating and registering of the functionality that is needed.


### Steps
You can start the database with only a csv.
Go over steps below, And see the result in your browser.
1. place csv file, in dir extras.
2. `python3 create_model.py > ../model.go`
3. cd ../
4. go fmt
5. go build
6. ./lambda --help
7. ./lambda  --csv assets/items.csv or `python3 ingestion.py -b 1000`
9. curl 127.0.0.1:8128/help/
10. browser 127.0.0.1:8128/


11. instructions curl 127.0.0.1:8128/help/ | python -m json.tool



### Running

sudo docker-compose up  --no-deps --build

promql {instance="lambdadb:8000"}

python3 extras/ingestion.py  -f movies_subset.tsv -format tsv -dbhost 127.0.0.1:8000
