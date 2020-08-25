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
docker-compose up
go get github.com/bitly/go-nsq
go get go.mongodb.org/mongo-driver/mongo
```
