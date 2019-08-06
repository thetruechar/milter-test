FROM registry.yixindev.net:5000/tls:latest
RUN ls /usr/local
ADD milter-test /usr/local/bin

RUN ls /usr/local/bin/milter
EXPOSE 9001

#ENTRYPOINT ["milter"]
#CMD ["service", "local"]

