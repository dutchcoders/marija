## Configuration
Go to the Marija installation directory and copy the sample configuration file.
``` 
cd go/src/github.com/dutchcoders/marija/
cp config.toml.sample config.toml 
```
If Elasticsearch is not installed, comment it out to prevent errors.
```
#[datasource]
#[datasource.elasticsearch]
#type="elasticsearch"
#url="http://127.0.0.1:9200/demo_index"
#username=
#password=

[datasource.twitter]
type="twitter"
consumer_key=""
consumer_secret=""
token=""
token_secret=""

[datasource.blockchain]
type="blockchain"

[[logging]]
output = "stdout"
level = "debug"
```