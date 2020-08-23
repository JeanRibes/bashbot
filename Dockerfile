FROM golang:1.14 as builder

WORKDIR /go/src
RUN mkdir -p bashbot

COPY go.mod bashbot
COPY go.sum bashbot

# téléchargement des dépendances
RUN cd bashbot && go mod tidy

COPY main.go bashbot

# build Go
#RUN CGO_ENABLED=0 GOOS=linux cd bashbot && go build -a -installsuffix cgo -o main main.go
RUN cd bashbot && go build  -o main main.go
# j'utilise ubuntu et pas alpine pour avoir plus de commande shell par défaut
RUN apt-get install ca-certificates && update-ca-certificates

RUN groupadd --gid 1000 bots && useradd -m --gid 1000 --uid 1000 bashbot
USER bashbot

ENV BOT_TOKEN "very secret"
CMD ["/go/src/bashbot/main"]