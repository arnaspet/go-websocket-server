FROM golang:1.11

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
SHELL ["/bin/bash", "-c"]

WORKDIR /go/src/teso_task
ENV PORT=3333

COPY .. .
RUN dep ensure
RUN go install -v ./...

CMD ["/bin/bash", "-c", "teso_task -port $PORT"]
