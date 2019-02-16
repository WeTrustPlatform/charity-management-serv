### Troubleshoot
- To connect to DB on EC2 instance:
```
sudo docker run --rm -it -v ~/:/workdir --env-file ./cms_env postgres:10-alpine /bin/sh
psql $DATABASE_URL
```
