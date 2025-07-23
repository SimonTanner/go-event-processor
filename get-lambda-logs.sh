CONTAINER="$(sudo docker container ls -a --format json | grep lambda | jq .Names)"
CONTAINER=`(echo ${CONTAINER//\"})`

sudo docker logs $CONTAINER