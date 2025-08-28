## Instructions
### MongoDB
#### Local Development
- To Bring up db infra:
```bash
$ cd ./infra/mongo
$ docker compose up --wait
```   
- To Bring down db infra   
```bash
$ cd ./infra/mongo
$ docker compose down
## The below command will also remove created volumes!!
$ docker compose down -v
```
#### Production
- Edit [mongo-admin.sh](mongo/mongo-admin.sh) and update configuration variables
- `bash mongo-admin.sh`
- Grab the password of `MONGO_ADMIN_USER` added in `mongo-admin.sh` for next steps
- Complete db migrations    
  ```bash
  $ docker run --rm --mount type=bind,source=./infra/mongo/migrations,target=/migrations \
       --network host migrate/migrate \
        -path=/migrations/ -database mongodb://mongo:{{ MONGO_ADMIN_PASS }}@localhost:27017/shop?authSource=admin up
  ```