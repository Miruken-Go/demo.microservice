FROM node:18.16.0-bookworm

#Install Docker
#https://docs.docker.com/engine/install/debian/
RUN apt-get update
RUN apt-get install ca-certificates curl gnupg

RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
RUN chmod a+r /etc/apt/keyrings/docker.gpg

RUN echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update
RUN VERSION_STRING=5:24.0.4-1~debian.12~bookworm; apt-get install -y docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io docker-buildx-plugin docker-compose-plugin

#Install AZ
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# The container build started failing because adding the containerapp extension needed pip
# Check and see if later versions can run without it
RUN apt-get -y install python3-pip 

# Need to install extensions in a system folder
# By default they do into the $HOME folder which is /root
# Irritatingly Github Actions overrides $HOME when running the container
RUN az extension add --name containerapp --system

CMD ["/bin/bash"]





