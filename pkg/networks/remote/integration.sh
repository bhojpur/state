#!/usr/bin/env bash

# Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

# XXX: this script is intended to be run from a fresh Digital Ocean droplet

# NOTE: you must set this manually now
echo "export DO_API_TOKEN=\"yourtoken\"" >> ~/.profile

sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install -y jq unzip python-pip software-properties-common make

# get and unpack golang
curl -O https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
tar -xvf go1.16.5.linux-amd64.tar.gz

## move binary and add to path
mv go /usr/local
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile

## create the goApps directory, set GOPATH, and put it on PATH
mkdir goApps
echo "export GOPATH=/root/goApps" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile
# **turn on the go module, default is auto. The value is off, if Bhojpur State source code
#is downloaded under $GOPATH/src directory
echo "export GO111MODULE=on" >> ~/.profile

source ~/.profile

mkdir -p $GOPATH/src/github.com/bhojpur
cd $GOPATH/src/github.com/bhojpur
# ** use git clone instead of go get.
# once go module is on, go get will download source code to
# specific version directory under $GOPATH/pkg/mod the make
# script will not work
git clone https://github.com/bhojpur/state.git
cd state
## build
make tools
make build
#** need to install the package, otherwise terdermint testnet will not execute
make install

# generate an ssh key
ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -N ''
echo "export SSH_KEY_FILE=\"\$HOME/.ssh/id_rsa.pub\"" >> ~/.profile
source ~/.profile

# install terraform
wget https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip
unzip terraform_0.11.7_linux_amd64.zip -d /usr/bin/

# install ansible
sudo apt-get update -y
sudo apt-add-repository ppa:ansible/ansible -y
sudo apt-get update -y
sudo apt-get install ansible -y

# required by ansible
pip install dopy

# the next two commands are directory sensitive
cd $GOPATH/src/github.com/bhojpur/state/pkg/networks/remote/terraform

terraform init
terraform apply -var DO_API_TOKEN="$DO_API_TOKEN" -var SSH_KEY_FILE="$SSH_KEY_FILE" -auto-approve

# let the droplets boot
sleep 60

# get the IPs
ip0=`terraform output -json public_ips | jq '.value[0]'`
ip1=`terraform output -json public_ips | jq '.value[1]'`
ip2=`terraform output -json public_ips | jq '.value[2]'`
ip3=`terraform output -json public_ips | jq '.value[3]'`

# to remove quotes
strip() {
  opt=$1
  temp="${opt%\"}"
  temp="${temp#\"}"
  echo $temp
}

ip0=$(strip $ip0)
ip1=$(strip $ip1)
ip2=$(strip $ip2)
ip3=$(strip $ip3)

# all the ansible commands are also directory specific
cd $GOPATH/src/github.com/bhojpur/state/pkg/networks/remote/ansible

# create config dirs
state testnet

ansible-playbook -i inventory/digital_ocean.py -l sentrynet install.yml
ansible-playbook -i inventory/digital_ocean.py -l sentrynet config.yml -e BINARY=$GOPATH/src/github.com/bhojpur/state/build/statectl -e CONFIGDIR=$GOPATH/src/github.com/bhojpur/state/pkg/networks/remote/ansible/mytestnet

sleep 10

# get each nodes ID then populate the ansible file
id0=`curl $ip0:26657/status | jq .result.node_info.id`
id1=`curl $ip1:26657/status | jq .result.node_info.id`
id2=`curl $ip2:26657/status | jq .result.node_info.id`
id3=`curl $ip3:26657/status | jq .result.node_info.id`

id0=$(strip $id0)
id1=$(strip $id1)
id2=$(strip $id2)
id3=$(strip $id3)

# remove file we'll re-write to with new info
old_ansible_file=$GOPATH/src/github.com/bhojpur/state/pkg/networks/remote/ansible/roles/install/templates/systemd.service.j2
rm $old_ansible_file

# need to populate the `--p2p.persistent-peers` flag
echo "[Unit]
Description={{service}}
Requires=network-online.target
After=network-online.target

[Service]
Restart=on-failure
User={{service}}
Group={{service}}
PermissionsStartOnly=true
ExecStart=/usr/bin/statectl node --mode validator --proxy-app=kvstore --p2p.persistent-peers=$id0@$ip0:26656,$id1@$ip1:26656,$id2@$ip2:26656,$id3@$ip3:26656
ExecReload=/bin/kill -HUP \$MAINPID
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
" >> $old_ansible_file

# now, we can re-run the install command
ansible-playbook -i inventory/digital_ocean.py -l sentrynet install.yml

# and finally restart it all
ansible-playbook -i inventory/digital_ocean.py -l sentrynet restart.yml

echo "congratulations, your testnet is now running :)"