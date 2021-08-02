# otter-cloud-ws
Setting config in config.ini.  

This project use minio as file sytem.  
You can learn how to use minio [here](https://docs.min.io/docs/minio-quickstart-guide.html).  
  
Demo: [https://www.calicomoomoo.com/otter-cloud/](https://www.calicomoomoo.com/otter-cloud/)   
Front End: [https://github.com/EricChiou/otter-cloud](https://github.com/EricChiou/otter-cloud)   
Test account:   
&nbsp;&nbsp;&nbsp;&nbsp;acc: test@gmail.com   
&nbsp;&nbsp;&nbsp;&nbsp;pwd: test   

### Making config file
Making config.ini with content below:  
```
# Server
SERVER_NAME=otter cloud
SERVER_PORT=7000
SSL_CERT_FILE_PATH=
SSL_KEY_FILE_PATH=

# MySQL
MYSQL_ADDR=127.0.0.1
MYSQL_PORT=3306
MYSQL_USERNAME=your user name
MYSQL_PASSWORD=your user password
MYSQL_DBNAME=your db name

# JWT
JWT_KEY=your jwt key
# JWT expire time, set 1 for one day, set 2 for two days, ...
JWT_EXPIRE=1

# RSA, key file path
RSA_PUBLIC_KEY=
RSA_PRIVATE_KEY=

# minio
END_POINT=
ACCESS_KEY_ID=
SECRET_ACCESS_KEY=
USE_SSL=
BUCKET_HASH_KEY=

# Environment
ENV=dev
```
