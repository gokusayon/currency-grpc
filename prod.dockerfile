# docker build --no-cache -f prod.dockerfile -t gokusayon/go-mongodb-micro .
# docker run -p 8000:9090 --name micro --network mongo-microservice_network1 gokusayon/go-mongodb-micro

# ------------- Builder Stage ---------------
FROM golang:1.14.3 as builder

# Setting up the env varaibles
ENV GO111MODULE=on
ARG SOURCE_LOCATION=/app

# Set up workdir
RUN mkdir -p ${SOURCE_LOCATION}
WORKDIR ${SOURCE_LOCATION}

# Downlaod dependencies
COPY go.mod .
RUN go mod download

# Build app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# ----------- Final Build ----------------
FROM alpine:latest  
ARG SOURCE_LOCATION=/app

# Install curl and set working dir
RUN apk --no-cache add curl
WORKDIR /root/

# Copy binaries from previous stage
COPY --from=builder ${SOURCE_LOCATION} .
EXPOSE 8080

# Docker container entrypoint
ENTRYPOINT ["./app"]