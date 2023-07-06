FROM node:18.16.0

#Install AZ
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN apt-get update

# The container build started failing because adding the containerapp extension needed pip
# Check and see if later versions can run without it
RUN apt-get -y install python3-pip 

# Need to install extensions in a system folder
# By default they do into the $HOME folder which is /root
# Irritatingly Github Actions overrides $HOME when running the container
RUN az extension add --name containerapp --system

CMD ["/bin/bash"]





