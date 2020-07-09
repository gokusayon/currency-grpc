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
CMD ["go", "run", "main.go"]