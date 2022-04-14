FROM golang:1.18 as godoc

WORKDIR /

RUN CGO_ENABLED=0 go install golang.org/x/tools/cmd/godoc@latest

COPY ./scripts/entrypoint.sh .

ENTRYPOINT ["./entrypoint.sh"]

