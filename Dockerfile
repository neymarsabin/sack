FROM golang:1.17.5
MAINTAINER neymarsabin <reddevil.sabin@gmail.com>

WORKDIR /app
COPY ./*.go ./

# build for linux
RUN CGO_ENABLED=0 go build -o /app/sack ./

# expose 6379 port
EXPOSE 6379

# run the binary
CMD ["/app/sack"]
