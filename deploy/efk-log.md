# The Necessity of the EFK System
During local deployments, logs are conveniently output to the `output/log/svc` directory, enabling easy log inspection. However, this approach proves inadequate when dealing with large data volumes or cloud - based services. Relying on the traditional method would require SSH access to the server for log queries. In a distributed deployment scenario, a single debugging session might demand connections to multiple servers, followed by manual retrieval of relevant logs from numerous files.

To address these challenges, we have implemented the **"EFK" (Elasticsearch - Filebeat - Kibana) system**. Filebeat, deployed on each server, collects local logs and forwards them to the Elasticsearch cluster. This setup allows for unified log queries on Kibana.

# Why EFK Instead of ELK?
Constrained by budget limitations and the simplicity of our requirements, we opted for the more lightweight Filebeat as a substitute for Logstash, which is a core component in the ELK stack.

# How to Query Logs?
Using the locally - deployed EFK as an example, we'll introduce some basic yet commonly used query methods in Kibana's Dev Tools. (Note: A more user - friendly visualization dashboard will be developed in the future for more intuitive log queries.)

## Query Steps:
1. Navigate to your Kibana homepage. If you haven't modified the Docker Compose configuration, the URL should be [kibana - home](http://localhost:5601/app/home#/).
2. Open the sidebar and locate the "Dev Tools" option under "Management" at the bottom. Click to enter the query interface.
3. Execute your queries.

If you're new to query statements, don't worry. The following are some simple examples. Just follow the provided comments and adjust the "size" and "from" parameters according to your needs. It's important to note that in Filebeat, we've set the log index name as `domtok - logs`.

### View the Structure of the logs index
```logstash
GET /domtok - logs
```

### Query All Docs within a Specified Range
```logstash
GET /domtok - logs/_search
{
  "query": {
    "match_all": {}
  },
  "from": 0, 
  "size": 20
}
```

### Query Docs of a Specified Service
```logstash
GET /domtok - logs/_search
{
  "query": {
    "match": {
      "service.keyword": "user"
    }
  },
  "from": 0, 
  "size": 20
}
```

### Query Docs of a Specified Service and a Specified Source. The default source format for a service is `app - serviceName`, e.g., `app - user`.
```logstash
GET /domtok - logs/_search
{
  "query": {
    "bool": {
      "must": [
        {"match": { "service.keyword": "user" }},
        {"match": { "source.keyword": "klog" }}
      ]
    }
  },
  "from": 0, 
  "size": 20
}
```

### Add a Match for the `msg` Content Based on the Previous Query
```logstash
GET /domtok - logs/_search
{
  "query": {
    "bool": {
      "must": [
        {"match": { "service.keyword": "user" }},
        {"match": { "source.keyword": "klog" }},
        {"match": { "msg": "etcd registry" }}
      ]
    }
  },
  "from": 0, 
  "size": 20
}
```

To send a query request, simply click the arrow on the right side of the query input area.
![img.png](img/kibana-dev-tools-sendRequest.png)
