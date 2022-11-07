FROM debian:10

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -sf /bin/true /sbin/initctl
ENV DEBIAN_FRONTEND noninteractive
ENV VIRTUAL_HOST localhost
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y apache2 php7.3 php7.3-bcmath php7.3-curl php7.3-gd \
	php7.3-imagick php7.3-imap php7.3-intl php7.3-json php7.3-ldap \
	php7.3-mbstring php7.3-mysql php7.3-opcache php7.3-readline \
	php7.3-redis php7.3-soap php7.3-xml php7.3-xmlrpc php7.3-xsl \
	php7.3-zip  libapache2-mod-php7.3  curl postfix vim  libcurl3-dev \
	telnet # openssh-server
ADD ./000-default.conf /etc/apache2/sites-available/000-default.conf
ADD ./php.ini /etc/php/7.3/apache2/php.ini
ADD ./moodle /var/www/moodle
ADD ./moodledata /var/www/moodledata_src
RUN chown -R www-data:www-data /var/www/moodle
EXPOSE 80
CMD ["apache2ctl", "-D", "FOREGROUND"]
