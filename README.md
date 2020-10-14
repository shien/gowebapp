# Golang で作る Web アプリケーション開発

動かしてみた

## 環境
* Lenovo Thinkpad X220
* Intel(R) Core(TM) i5-2520M CPU @ 2.50GHz
* Ubuntu 20.04 LTS
* go version go1.14 linux/amd64 

## 作者によるサンプルコード

https://github.com/matryer/goblueprints

## 環境



## Chat application


BootStrap css のバージョンだけ 4.5.0 にあげてある

```
cp -r trace/ ${GOPATH}/src/
```

main.go のsecret key Facebook / Google / GitHub のどれかを入力

```
sh buildrun.sh
```

http://localhost:8080/chat

にアクセスする

## CLIs


### sprinkle

```
go build -o sprinkle
```

### coolify

```
go build -o coolify
```
### domainify

```
go build -o domainify
```

### synonyms

https://words.bighugelabs.com/

でアカウントを作成し、API KEY を取得

```
export BTH_APIKEY=[API_KEY]
cp -r thesaurus ${GOPATH}/src
go build -o synonyms
```

### domainfinder

```
cd domainfinder
sh build.sh
./domainfinder
```

## Distributed System and flexible data processing

```
go get github.com/bitly/go-nsq
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joeshaw/envdecode
go get github.com/gomodule/oauth1/oauth
go get "go.mongodb.org/mongo-driver/bson"
docker-compose up
go build -o twittervotes
./twittervotes
```

## REST形式でデータや機能を公開

```
go get github.com/stretchr/graceful
```

API サーバー は `-addr` 無指定で localhost:8080 で起動
```
docker-compose up
go build -o api
./api -mongo 172.20.0.12:27017
```

Web サーバーに関しては  jsapi が終了したみたいなので写経だけ…

https://developers-jp.googleblog.com/2016/07/google-feed-api.html
