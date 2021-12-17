FROM golang:alpine as build_go_bin
WORKDIR /go/src/videochat
ADD . .
RUN ls -lag
RUN go build -o cv_server

FROM alpine:latest
COPY --from=build_go_bin /go/src/videochat/cv_server /app/cv_server
COPY --from=build_go_bin /go/src/videochat/public /public
CMD [ "/app/cv_server" ]
