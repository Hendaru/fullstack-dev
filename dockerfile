FROM golang:1.19.2-alpine3.15

WORKDIR /go/src/github.com/Hendaru/bwaStartup


COPY . .

# RUN go build -o bwa-backend
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build main.go" -command="./main"

# EXPOSE 8080

# CMD ./bwa-backend

# FROM golang:1.16-alpine
# #FROM golang:1.19.2-alpine3.15
# WORKDIR /app


# # add some necessary packages
# RUN apk update && \
#     apk add libc-dev && \
#     apk add gcc && \
#     apk add make



# # prevent the re-installation of vendors at every change in the source code
# COPY ./go.mod go.sum ./
# RUN go mod download && go mod verify

# # # Install Compile Daemon for go. We'll use it to watch changes in go files
# RUN go get github.com/githubnemo/CompileDaemon

# # # Copy and build the app
# COPY . .
# COPY ./entrypoint.sh /entrypoint.sh

# # # wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
# ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
# RUN chmod +rx /usr/local/bin/wait-for ./entrypoint.sh

# ENTRYPOINT [ "sh", "/entrypoint.sh"]
