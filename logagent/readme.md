
####启动kafka消息客户端
`bin\windows\kafka-console-consumer.bat --bootstrap-server localhost:9092 --topic web_log --from-beginning`


#### grafana连接influxdb,grafana登录账号及密码，admin,admin
- 在URL栏输入http://localhost:8086,默认显示的是placeholder,
- 在Custom HTTP Headers添加一个key:Authorization,value为Token token的值
- eg`hedaer:Authorization`  `value:Token sdhfogk342klndsfdsdsbfdsfdsfkdsfdfd-sdf=sdfsddf`,Token和值之间有个空格


####influx
- 列出bucket,-o表示组名
`influx bucket list -o loocc -t P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ==`
  
- 新版需要把Bucket映射成database，不然grafana中找不到数据库,在用api操作时注意Bucket名和measurement名一致
- 第一在grafana连接influxdb时，会指定数据库，如果后面需要查询其它库中的表，需要修改连接中的database或者重新建一个连接

  
  ```
  curl --request POST http://localhost:8086/api/v2/dbrps \
  --header "Authorization: Token P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ==" \
  --header 'Content-type: application/json' \
  --data '{
        "bucketID": "13270780354c09e3",
        "database": "monitor",
        "default": true,
        "orgID": "27929176501354b0",
        "retention_policy": "120d"
      }'
```

curl --get http://localhost:8086/query?db=monitor \
  --header "Authorization: Token P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ==" \
  --data-urlencode "q=SELECT cpu_percent FROM monitor"


curl --get http://localhost:8086/query?db=test \
  --header "Authorization: Token P5AnFpfNP3DVbd-ZH-S40HGg5b91Xsj1FCOCEpX-kcV28t7-rPJGrzX3riEQl8RijfND-_djFyNJtqNo8p4QwQ==" \
  --data-urlencode "q=SELECT max FROM test"
