# LambdaDB
In memory database that uses filters to get the data you need.
Lambda DB has a tiny codebase which does a lot
Lambda is not ment as a persistance storage or a replacement for a traditional
Database but as fast analytics engine cache representation engine.

powers: https://dego.vng.nl

## Properties:

- Insanely fast API. 1ms respsonses
- Fast to setup.
- Easy to deploy.
- Easy to customize.
- Easy export data

- Implement custom authorized filters.

## Indexes

- S2 geoindex for fast point lookup
- Bitarrays
- Mapping

- Your own special needs indexes!

## Flow:

Generate a model and load your data.
The API is generated from your model.
Deploy.

Condition: Your dataset must fit in memory.

Can be used for your needs by changing the `models.go` file to your needs.
Creating and registering of the functionality that is needed.


### Steps
You can start the database with only a csv.
Go over steps below, And see the result in your browser.

1. place csv file, in dir extras.
2. `python3 create_model_.py`  answer the questions.
3. go fmt model.go
4. mv model.go ../
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
=======

1. instructions curl 127.0.0.1:8000/help/ | python -m json.tool

### Questions



### TODO

- load data directly from a database (periodic)
- document the `create_model.py` questions
- use a remote data source
- use some more efficient storage method (done)
- generate swagger API
- Add more tests
