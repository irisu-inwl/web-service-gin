FROM golang as builder

RUN apt update && apt install git
RUN mkdir /opt/myapp
WORKDIR /opt/myapp
COPY ./src .
CMD go get .

FROM builder as runner

CMD go run .

FROM builder as test

CMD go test ./...
