FROM moby/buildkit:v0.9.3
WORKDIR /state
COPY state README.md /state/
ENV PATH=/state:$PATH
ENTRYPOINT [ "/bhojpur/state" ]