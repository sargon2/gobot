FROM alpine:3.11.3

ARG USER_ID
ARG GROUP_ID

RUN addgroup --gid $GROUP_ID user
RUN adduser --disabled-password --gecos '' --uid $USER_ID --ingroup user user

RUN apk add go=1.13.13-r0

RUN mkdir /gobot
RUN chown user:user /gobot
RUN mkdir /opt/mount # Used for copying build results back out
RUN chown user:user /opt/mount

WORKDIR /gobot

USER user

COPY --chown=user:user go.mod go.sum ./

RUN go mod download

# Copy everything except go.mod and go.sum in
RUN mkdir /tmp/backup && cp go.mod go.sum /tmp/backup/
COPY --chown=user:user . .
RUN cp -f /tmp/backup/* . && rm -rf /tmp/backup/

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
