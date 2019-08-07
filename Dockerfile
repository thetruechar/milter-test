FROM registry.yixindev.net:5000/tls:latest
ADD milter-test /usr/local/bin

EXPOSE 9001

#ENTRYPOINT ["milter"]
#CMD ["service", "local"]

