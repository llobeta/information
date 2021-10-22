FROM centos
RUN yum update -y
RUN yum install wget -y
RUN wget https://downloads.mariadb.com/MariaDB/mariadb_repo_setup
RUN echo "fc84b8954141ed3c59ac7a1adfc8051c93171bae7ba34d7f9aeecd3b148f1527 mariadb_repo_setup" \
    | sha256sum -c -
RUN chmod +x mariadb_repo_setup
RUN ./mariadb_repo_setup \
   --mariadb-server-version="mariadb-10.5"
RUN yum install MariaDB-client -y
COPY InformationService/InformationService .
CMD ./InformationService
