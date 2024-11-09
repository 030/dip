FROM golang:1.23.3-alpine3.17 as builder
ENV USERNAME=dip
ENV BASE=/opt/${USERNAME}
COPY . /go/${USERNAME}/
WORKDIR /go/${USERNAME}/cmd/${USERNAME}/
RUN apk add --no-cache curl=~7 git=~2 && \
    CGO_ENABLED=0 go build && \
    curl -sL https://gist.githubusercontent.com/030/54fc7ae735a163c09dcf6f3699d87e81/raw/b82514f50525e0ebf843e0dbf9bef1a382ccd40f/openshift-docker-user-entrypoint.sh > entrypoint.sh && \
    curl -sL https://gist.githubusercontent.com/030/34a2bf3f7f1cd427dc36c86dcb1e8cf7/raw/e1be7ef3c2c1a8441e406a669a0f6b6d97dcc984/openshift-docker-user.sh > user.sh && \
    chmod +x user.sh && \
    ./user.sh

FROM alpine:3.20.1
ENV BIN=/usr/local/bin/
ENV USERNAME=dip
ENV BASE=/opt/${USERNAME}
ENV BASE_BIN=${BASE}/bin
ENV PATH=${BASE_BIN}:${PATH}
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /opt/ /opt/
USER $USERNAME
ENTRYPOINT ["entrypoint.sh"]
CMD ["dip"]
