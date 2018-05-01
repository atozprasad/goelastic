export APP_ELASTIC_INDEX_NAME="products"
export APP_ELASTIC_TYPE_NAME="product"
export APP_ELASTIC_URL1="http://35.202.99.46:9200"
export APP_ELASTIC_URL2="http://35.192.32.150:9200"
export APP_ELASTIC_URL3="http://35.224.21.162:9200"
export APP_ELASTIC_USER="elastic"
export APP_ELASTIC_PASSWORD="hKVd9xXQ"
export GIN_MODE=debug

hostos=`uname -a`
./autocomplete_0.1_linux_amd64

if [ $hostos == "Linux" ]
then
	./autocomplete_0.1_linux_amd64 &
else
	./autocomplete_0.1_darwin_amd64 &
fi
