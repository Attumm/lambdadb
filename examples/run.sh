# run from root of the project
# bash example/run.sh

python3 extras/create_model.py -f examples/movies_subset.tsv -format tsv > model.go
go build -o lambda_db
go fmt

bash examples/ingest.sh &

./lambda_db -indexed y -frontend y


