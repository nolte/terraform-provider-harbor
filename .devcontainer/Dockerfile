FROM docker.pkg.github.com/nolte/vscode-devcontainers/k8s-operator:latest



# https://github.com/bats-core/bats-core
ENV BATS_VERSION=1.2.0

ENV SWAGGER_VERSION=v0.23.0

USER root

COPY --from=docker.pkg.github.com/nolte/vscode-devcontainers/devops:latest /usr/local/bin/terraform /usr/local/bin/terraform
RUN curl -sSL -k https://github.com/bats-core/bats-core/archive/v${BATS_VERSION}.tar.gz -o /tmp/bats.tgz \
  && tar -zxf /tmp/bats.tgz -C /tmp \
  && /bin/bash /tmp/bats-core-${BATS_VERSION}/install.sh /usr/local

RUN curl -Lo /usr/local/bin/swagger https://github.com/go-swagger/go-swagger/releases/download/${SWAGGER_VERSION}/swagger_linux_amd64 \
  && chmod +x /usr/local/bin/swagger 

ENV TERRAFORM_DOC_VERSION=v0.9.1

RUN curl -Lo /usr/local/bin/terraform-doc https://github.com/segmentio/terraform-docs/releases/download/${TERRAFORM_DOC_VERSION}/terraform-docs-${TERRAFORM_DOC_VERSION}-$(uname | tr '[:upper:]' '[:lower:]')-amd64 \
  && chmod +x /usr/local/bin/terraform-doc 


RUN curl -sSL -k https://github.com/bflad/tfproviderdocs/releases/download/v0.6.0/tfproviderdocs_0.6.0_linux_amd64.tar.gz -o /tmp/tfproviderdocs.tgz \
  && tar -zxf /tmp/tfproviderdocs.tgz -C /tmp \
  && mv /tmp/tfproviderdocs /usr/local/bin/


RUN apk add --update-cache \
  nodejs npm \
  && rm -rf /var/cache/apk/* 


RUN mkdir -p /home/$USERNAME/.vscode-server/extensions \
  /home/$USERNAME/.vscode-server-insiders/extensions \
  && chown -R $USERNAME \
  /home/$USERNAME/.vscode-server \
  /home/$USERNAME/.vscode-server-insiders

RUN mkdir -p /go/src && chown -R ${USER_UID}:${USER_GID} /go/src \
  && mkdir -p /go/pkg && chown -R ${USER_UID}:${USER_GID} /go/pkg

USER ${USERNAME}

#COPY --chown=${USER_UID}:${USER_GID} --from=docker.pkg.github.com/nolte/vscode-devcontainers/devops:latest  /home/${USERNAME}/.zshrc-specific /home/${USERNAME}/.zshrc-specific
COPY --chown=${USER_UID}:${USER_GID} files/devops-zshrc-specific /home/${USERNAME}/.zshrc-specific
COPY --chown=${USER_UID}:${USER_GID} files/devops-welcome.sh /home/${USERNAME}/.welcome.sh

RUN mkdir -p /home/${USERNAME}/.terraform.d/plugins/linux_amd64

ENV PATH="/usr/local/share/npm/bin:/usr/local/bin:/usr/local/sbin:~/bin:$PATH"

RUN mkdir "/home/${USERNAME}/.npm-packages"
RUN npm config set prefix "/home/${USERNAME}/.npm-packages"
RUN npm install swagger-merger --user -g

RUN helm repo add harbor https://helm.goharbor.io
RUN npm install swagger-merger -g