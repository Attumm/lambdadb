# LambdaDB
In memory database that uses filters to get the data you need.

Can be used for your needs by changing the models.go file to your needs.
Creating and registering of the functionality that is needed.


### example
LambdaDB loaded with dataset from imdb at around 7 million items.
Frontend of LambdaDB shows the database in action.
![LambdaDB](https://imgur.com/JPGAb3w)

### Steps
You can start the database with only a csv.
Go over steps below, And see the result in your browser.
1. `python3 extras/create_model.py -f <path_to_file> ../model.go`
2. go fmt
3. go build
4. ./lambdadb --help
5. python3 extras/ingestion.py -f  <path_to_file>
6. curl 127.0.0.1:8128/help/
7. browser http://127.0.0.1:8128/
8. examples curl 127.0.0.1:8128/help/ | python3 -m json.tool

### Create Snapshot
http://127.0.0.1:8128/mgmt/save
 
### Load Snapshot
http://127.0.0.1:8128/mgmt/load
 
### Use index
Currently the index is on all the columns.
To run the index start lambdadb with indexed.
Create a snapshot of the current data compressed.
1. `http://127.0.0.1:8128/mgmt/save/bytesz`
2. `./lambda_db -indexed`
3. `http://127.0.0.1:8128/mgmt/load/bytesz`

### Running

sudo docker-compose up  --no-deps --build

promql {instance="lambdadb:8000"}

python3 extras/ingestion.py  -f movies_subset.tsv -format tsv -dbhost 127.0.0.1:8000
