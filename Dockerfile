FROM ubuntu:14.04

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates

COPY ./bin/selfie-linux-amd64 /usr/bin/selfie
COPY ./etc/selfie-prod.conf /etc/selfie.conf
RUN mkdir /usr/bin/temp
RUN mkdir /usr/bin/bundle

EXPOSE 7331

CMD ["selfie", "-config=/etc/selfie.conf"]
