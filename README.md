# Run

go run main.go data_type.go util.go db.go

# using docker

sudo docker run --name=crawler_go --mount source=output,destination=/app/output --mount source=db,destination=/app/db linhabc/nhattao_crawler

sudo docker run -d -it --name crawler_go --mount type=bind,source=/home/linhnguyen/app_target/db,target=/app/db --mount type=bind,source=/home/linhnguyen/app_target/output,target=/app/output linhabc/nhattao_crawler

sudo docker run -d -it --name crawler_go --mount type=bind,source=/u01/data/nhattao.net/db,target=/app/db --mount type=bind,source=/u01/data/nhattao.net/output,target=/app/output linhabc/nhattao_crawler.....

# Using docker-compose

docker-compose up

# output folder

- output: store generated json file
- db: store generated leveldb folder
