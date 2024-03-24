# TinyURL

A higly available url shortener written in golang

## Functional Requirments
User need to be redirected from the shortened url to the actual one
## Non-Functional Requirments
- The application needs to able to handle about 100k requests per day without going down
- Need to be self-healing, if it goes down it need to be able to start again fastly
- Needs to be 12 factor complient
- Needs to be able to run relatively cheap
- Needs to be easily deployable
- Needs to be monitorable with dashboards etc

