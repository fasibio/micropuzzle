FROM gcr.io/distroless/base-debian10
ARG commit_sha
ARG application_build_id
ENV COMMIT_SHA=${commit_sha}
ENV APPLICATION_BUILD_ID=${application_build_id}
COPY micropuzzle /
CMD ["/micropuzzle"]
