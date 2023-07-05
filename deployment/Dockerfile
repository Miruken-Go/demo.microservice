FROM node:18.16.0

#Install AZ
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash
RUN apt-get update

# The container build started failing because adding the containerapp extension needed pip
# Check and see if later versions can run without it
RUN apt-get -y install python3-pip 

RUN az extension add --name containerapp --upgrade

CMD ["/bin/bash"]





