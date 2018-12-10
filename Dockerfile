FROM ubuntu:xenial

MAINTAINER http://ydk.io

COPY . /root/ydk-gen

RUN echo 'Installing dependencies'

WORKDIR /root/ydk-gen

RUN /bin/bash -c './test/dependencies_ubuntu.sh && ./test/dependencies_linux_gnmi.sh && pip install -r requirements.txt && ./generate.py --libydk && make -C gen-api/cpp/ydk/build install && ./generate.py --cpp --service profiles/services/gnmi-0.4.0.json && make -C gen-api/cpp/ydk-service-gnmi/build install && ./generate.py --core --python && pip install gen-api/python/ydk/dist/ydk*.tar.gz && ./generate.py --python --service profiles/services/gnmi-0.4.0.json && pip install gen-api/python/ydk-service-gnmi/dist/ydk*.tar.gz'
