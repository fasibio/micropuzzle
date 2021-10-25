FROM gcr.io/distroless/base-debian10
ARG buildNumber
ENV BUILDNUMBER=${buildNumber} 
COPY micropuzzle /
CMD ["/micropuzzle"]
