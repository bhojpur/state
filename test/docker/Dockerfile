FROM golang:1.17

# Grab deps (jq, hexdump, xxd, killall)
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  jq bsdmainutils vim-common psmisc netcat curl

# Setup Bhojpur State repo
ENV REPO $GOPATH/src/github.com/bhojpur/state
ENV GOBIN $GOPATH/bin
WORKDIR $REPO

# Copy in the code
# TODO: rewrite to only copy Makefile & other files?
COPY . $REPO

# Install the vendored dependencies
# docker caching prevents reinstall on code change!
RUN make tools

# install ABCI CLI
RUN make install_abci

# install Bhojpur State
RUN make install

RUN statectl testnet \
  --config $REPO/test/docker/config-template.toml \
  --node-dir-prefix="mach" \
  --v=4 \
  --populate-persistent-peers=false \
  --o=$REPO/test/p2p/data

# Now copy in the code
# NOTE: this will overwrite whatever is in vendor/
COPY . $REPO

# expose the volume for debugging
VOLUME $REPO

EXPOSE 26656
EXPOSE 26657