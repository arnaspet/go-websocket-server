FROM golang:1.11 as builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/arnaspet/teso_task/subscriber

COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/subscriber .


FROM scratch
COPY --from=builder /bin/subscriber /bin/subscriber
ENTRYPOINT ["/bin/subscriber"]
