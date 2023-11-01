FROM node:18.16.0-bookworm

############
#go
############
RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
ENV PATH "$PATH:/usr/local/go/bin"

############
#Docker
############
#https://docs.docker.com/engine/install/debian/
RUN apt-get update
RUN apt-get -y install ca-certificates curl gnupg

RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
RUN chmod a+r /etc/apt/keyrings/docker.gpg

RUN echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update
RUN VERSION_STRING=5:24.0.4-1~debian.12~bookworm; apt-get install -y docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io docker-buildx-plugin docker-compose-plugin

############
#AZ
############
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# The container build started failing because adding the containerapp extension needed pip
# Check and see if later versions can run without it
RUN apt-get -y install python3-pip 

# Need to install extensions in a system folder
# By default they do into the $HOME folder which is /root
# Irritatingly Github Actions overrides $HOME when running the container
RUN az extension add --name containerapp --system

############
#gh
############
RUN type -p curl >/dev/null || ( apt update &&  apt install curl -y)
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg |  dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
RUN chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" |  tee /etc/apt/sources.list.d/github-cli.list > /dev/null
RUN apt update
RUN apt install gh -y

CMD ["/bin/bash"]
