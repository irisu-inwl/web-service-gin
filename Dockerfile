FROM golang

RUN apt update && apt install git
RUN mkdir /opt/myapp
WORKDIR /opt/myapp
ADD ./src .

CMD bash