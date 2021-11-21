FROM golang:latest

RUN go get github.com/gin-gonic/gin@v1.7.1 && \
  go get github.com/lib/pq

WORKDIR /Users/monster/Desktop/modanisa/back-end

COPY . /Users/monster/Desktop/modanisa/back-end


EXPOSE 8080
CMD ["go", "run", "main.go"]