# cms


## 安装中间件

### mysql

1. sudo docker volume create mysqlData
2. sudo docker volume create mysqlConfig
3. ```
   docker run -d \
   --name mysql \
   --restart unless-stopped \
   -v mysqlData:/var/lib/mysql \
   -v mysqlConfig:/etc/mysql/conf.d \
   -e MYSQL_ROOT_PASSWORD=gateway#GWS2025 \
   -p 3306:3306 \
   mysql
   ```
4. sudo vim /var/lib/docker/volumes/mysqlConfig/_data/my.cnf
    ```
    [mysqld]
    #允许所有ip访问
    bind-address = 0.0.0.0
    character-set-server = utf8mb4
    collation-server = utf8mb4_unicode_ci
    # 调整连接数和缓存
    max_connections = 200
    
    # 时区配置
    default-time-zone = '+8:00'
    
    # 日志配置
    general_log = 1
    general_log_file = /var/log/mysql/general.log
    
    [client]
    default-character-set = utf8mb4
    
    [mysql]
    default-character-set = utf8mb4
    ```