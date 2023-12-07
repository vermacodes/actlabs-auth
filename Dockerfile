FROM actlab.azurecr.io/repro_base

WORKDIR /app

ADD actlabs-auth ./

EXPOSE 80/tcp
EXPOSE 443/tcp

ENTRYPOINT [ "/bin/bash", "-c", "export PORT='80' && ./actlabs-auth" ]