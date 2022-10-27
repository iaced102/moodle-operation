FROM debian:10

# ENV TZ=Asia/Ho_Chi_Minh
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -sf /bin/true /sbin/initctl
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y apache2 php7.3 php7.3-bcmath php7.3-curl php7.3-gd \
		php7.3-imagick php7.3-imap php7.3-intl php7.3-json php7.3-ldap \
		php7.3-mbstring php7.3-mysql php7.3-opcache php7.3-readline \
		php7.3-redis php7.3-soap php7.3-xml php7.3-xmlrpc php7.3-xsl \
		php7.3-zip  libapache2-mod-php7.3  curl postfix vim  libcurl3-dev
