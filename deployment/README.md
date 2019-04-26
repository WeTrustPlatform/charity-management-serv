### Troubleshoot
- To connect to DB on EC2 instance:
```
sudo docker run --rm -it -v ~/:/workdir --env-file ./cms_env postgres:10-alpine /bin/sh
psql $DATABASE_URL
```


- Manually seed on development instance:
```
sudo docker run -it --rm -d --network backend-net --env-file ./.cms_env -v ~/seed:/seed sihoang/charity-management-serv cms-seeder -data /seed/data-download-pub78.txt
```
Notes: Make sure to download data to `~/seed` folder in advance
```
wget -P ~/seed https://apps.irs.gov/pub/epostcard/data-download-pub78.zip
```


- The re-index takes a while. So restarting server or seeding will take a while to load DB. Be patient.
