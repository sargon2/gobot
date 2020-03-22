FROM alpine:3.11.3

ARG USER_ID
ARG GROUP_ID

RUN addgroup --gid $GROUP_ID user
RUN adduser --disabled-password --gecos '' --uid $USER_ID --ingroup user user

RUN apk add go=1.13.4-r1

WORKDIR /gobot

USER user

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test ./...
