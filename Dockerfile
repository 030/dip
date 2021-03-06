FROM golang:1.16.6-alpine3.14 as builder
ENV PROJECT dip
RUN mkdir $PROJECT && \
    adduser -D -g '' $PROJECT
COPY cmd ./$PROJECT/cmd/
COPY internal ./$PROJECT/internal/
COPY pkg ./$PROJECT/pkg/
COPY go.mod go.sum ./$PROJECT/
WORKDIR $PROJECT/cmd/$PROJECT
RUN apk add git && \
    CGO_ENABLED=0 go build && \
    cp $PROJECT /$PROJECT

FROM alpine:3.13.5
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /dip /usr/local/bin/dip
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER dip
ENTRYPOINT ["dip"]
