FROM cgr.dev/chainguard/go:latest AS gobuilder
# use static link build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/local/src/repo
COPY . /usr/local/src/repo
RUN go build ./cmd/switchboard

FROM cgr.dev/chainguard/wolfi-base
COPY --from=gobuilder /usr/local/src/repo/switchboard /bin/switchboard

ENTRYPOINT ["/bin/switchboard"]
