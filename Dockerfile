FROM scratch
MAINTAINER DeedleFake

COPY yabs /

VOLUME /etc/yabs/

ENTRYPOINT ["/yabs"]
