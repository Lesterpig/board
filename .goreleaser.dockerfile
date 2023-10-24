FROM scratch
# goreleaser already has the binaries, so we have a separate Dockerfile for this

COPY --chown=1000:1000 tcs-board /tcs-board

USER 1000:1000

ENTRYPOINT ["/tcs-board"]
