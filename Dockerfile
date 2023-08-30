FROM golang:alpine

RUN go version
ENV GOPATH=/

RUN apk add --no-cache make

COPY ./ ./
RUN go mod download
RUN make build

CMD [ "./bin/segments" ]