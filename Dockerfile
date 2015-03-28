FROM golang:1.4.1

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/GreatestGuys/pifuxelck-server-go

RUN go-wrapper download github.com/GreatestGuys/pifuxelck-server-go
RUN go-wrapper install github.com/GreatestGuys/pifuxelck-server-go

ENTRYPOINT ["/go/bin/pifuxelck-server-go"]
CMD ["--help"]
EXPOSE 3000
