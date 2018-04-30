---
api_version: go1
runtime: go
env_variables:
  APP_ELASTIC_INDEX_NAME: products
  APP_ELASTIC_PASSWORD: hKVd9xXQ
  APP_ELASTIC_TYPE_NAME: product
  APP_ELASTIC_URL1: "http://35.202.99.46:9200"
  APP_ELASTIC_URL2: "http://35.192.32.150:9200"
  APP_ELASTIC_URL3: "http://35.224.21.162:9200"
  APP_ELASTIC_USER: elastic
  GIN_MODE: debug
handlers:
- url: /.*
  script: _go_app

env: flex

# This sample incurs costs to run on the App Engine flexible environment.
# The settings below are to reduce costs during testing and are not appropriate
# for production use. For more information, see:
# https://cloud.google.com/appengine/docs/flexible/python/configuring-your-app-with-app-yaml
resources:
  cpu: 1
  memory_gb: 4
  disk_size_gb: 10
network:
  forwarded_ports:
    - 8080
    - 8080:8080
    - 8080/tcp
    - 80/tcp
