FROM alpine:3.11.3

ARG USER_ID
ARG GROUP_ID

RUN addgroup --gid $GROUP_ID user
RUN adduser --disabled-password --gecos '' --uid $USER_ID --ingroup user user

RUN apk add go=1.13.4-r1

RUN mkdir /gobot
RUN chown user:user /gobot

WORKDIR /gobot

USER user

COPY --chown=user:user go.mod go.sum ./

RUN go mod download

# Copy everything except go.mod and go.sum in
RUN mkdir /tmp/backup && cp go.mod go.sum /tmp/backup/
COPY --chown=user:user . .
RUN cp -f /tmp/backup/* . && rm -rf /tmp/backup/

RUN go mod tidy

RUN go test ./...
