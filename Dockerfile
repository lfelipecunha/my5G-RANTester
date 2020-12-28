FROM golang:1.14.4-stretch

RUN apt-get update && apt -y install build-essential libgsl-dev

RUN go get github.com/derekparker/delve/cmd/dlv

WORKDIR /usr/local/go/src/my5G-RANTester

COPY go.mod .
RUN go mod download

VOLUME ["/usr/local/go/src/my5G-RANTester/internal"]
VOLUME ["/usr/local/go/src/my5G-RANTester/cmd"]
VOLUME ["/usr/local/go/src/my5G-RANTester/config"]
VOLUME ["/usr/local/go/src/my5G-RANTester/lib"]

#RUN cd src && go build -o my5g-rantester

#ENTRYPOINT ["/my5G-RANTester/src/my5g-rantester"]
CMD [ "dlv", "debug", "my5G-RANTester/cmd", "--listen=:40000", "--headless=true", "--api-version=2", "--log" , "--", "load-test" ]
