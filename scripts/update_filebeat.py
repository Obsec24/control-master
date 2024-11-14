import os
f = open("/etc/filebeat/filebeat.yml","r")
lineas = f.readlines()
f.close()
f = open("/etc/filebeat/filebeat.yml","w")
elasticsearch_password = os.getenv('ELASTICSEARCH_PASSWORD')
elasticsearch_host =os.getenv('ELASTICSEARCH_HOST')
for l in lineas:
    if l == '  password: "changeme"'+"\n":
        l = '  password: ' +'"'+ elasticsearch_password + '"'+ "\n"
    if l == '  hosts: ["192.168.1.201:9200"]'+"\n":
        l = '  hosts: '+'["'+elasticsearch_host + '"]'+ "\n"
    f.write(l)
f.close()
