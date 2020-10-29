# accounts-statistics-tool

A tool to fetch accounts data

### Prerequisites
- [Golang 1.12+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [MongoDB](https://www.mongodb.com/try/download)

### Environment variables
The following environment variables are used for program execution

Variable                      |Required  |Example                    |Default |Notes
------------------------------|----------|---------------------------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------
TRANSACTIONS_MONGODB_URL      | &#x2713; | mongodb://localhost:27017/|        | This variable must follow the standardised [MongoDB connection string format](https://docs.mongodb.com/manual/reference/connection-string/)
TRANSACTIONS_MONGODB_DATABASE | &#x2713; | users_application         |        |
LOG_LEVEL                     | &#x2717; | debug                     | info   | A lower case representation of the standard log level enumerations. Possible values can be found [here](https://github.com/sirupsen/logrus/blob/master/logrus.go#L25)
SENDER_EMAIL                  | &#x2713; | example@provider.co.uk    |        | This variable must be on the AWS authorised list
RECEIVER_EMAIL                | &#x2713; | example@provider.co.uk    |        | This variable must be on the AWS authorised list
SES_AWS_REGION                | &#x2713; | eu-west-1                 |        | This variable must be a valid AWS region

### Building and running

#### Docker
Bake a Docker image using the following command at the base of the project:

```
docker build --build-arg transactions_mongodb_url="<your_mongo_connection_string>" --build-arg transactions_mongodb_database="<your_mongo_db_name> --build-arg sender_email="<example@provider.co.uk>" --build-arg receiver_email="<example@provider.co.uk>" --build-arg ses_aws_region="<ses_aws_region_here>" -t <image_name> .
```

An additional `--build-arg` flag with key log_level can be optionally set to configure the log level of the environment.

Once built, run the image using:

```
docker run <image_name>
```
