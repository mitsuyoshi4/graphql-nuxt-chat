FROM golang:1.11

COPY . /go/src/github.com/mitsuyoshi4/graphqlchat

WORKDIR /go/src/github.com/mitsuyoshi4/graphqlchat

RUN make setup && make build

CMD ["./bin/GraphQLChat"]
