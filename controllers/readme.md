## MySQL database user creation process：**
    1、CREATE USER 'username'@'host' IDENTIFIED BY 'password';  
    2、GRANT ALL ON *.* TO 'username'@'host';  
    3、flush privileges;
## mysql login command：
    mysql -u root -p 