# ga-hsl-hrt

Simple GO based Google Assistant Action to retrieve Helsinki Regional Transport (HSL HRT) routes to pre-defined destinations.
Use case: Find bus timings from bus stops near to a house. E.g.
"Hey Google: When is the next 215 to Sello?"
"Ok Google: When is the next bus to Tapiola?"

## Getting Started

1. Clone the project
   
2. Update config-file.json with following configuration parameter for the application:
   
   1. Update interested bus route numbers against "routes".
   
   2. Update "callSignToHeadSign" with interested destinations.
      1. Finnish words are difficult to comprehend in Google Assistant conversations.
      2. Many destination names are simply too long and difficult to get through to Google Assistant.
      3. Hence this structure maps the actual destination names with simple invokable keywords.
      4. For e.g. in the supplied config-file.json, actual destination is "Leppävaara" and nickname is "Sello". End-users use Sello in conversations to find routes to Leppävaara.
   
   3. Update "stopGtfsIds". These are unique IDs for each stop and can be found this way:
      1. Identify the bus stop id using google maps or https://reittiopas.hsl.fi/.
         1. E.g. Search for Jupperinympyrä and it shows stop id as E1439.
         2. Open link - https://api.digitransit.fi/graphiql/hsl in a browser and type following query on left side:
            {
                stops(name: "E1438") {
                    gtfsId
                    name
                    code
                }
            }
        
            Press Play. This produces the details on the right side:
            {
                "data": {
                    "stops": [
                    {
                        "gtfsId": "HSL:2143218",
                        "name": "Jupperinympyrä",
                        "code": "E1439"
                    }
                    ]
                }
            }
            gtfsId of the stop with code E1439 is "HSL:2143218".
   
   4. Update server listening port under "port". Ensure this port is free, since this is the port the application will listen to and Google Assistant will try to access when invoking the action
   
   5. Update server TLS certificate location against "serverCert".
   
   6. Update server encryption key location against "serverKey".
   
   7. Update client certificate location against "clientCert". This is needed for mutual TLS.
   
   8. Update application log file location.

### Prerequisites

1. A working GO environment. Follow installation instructions from here - https://golang.org/dl/
   
2. Install other required GO packages. E.g. in a ubuntu shell:
   1.  go get -v github.com/spf13/viper
   2.  go get -v github.com/sirupsen/logrus
   3.  go get -v github.com/machinebox/graphql
   4.  go get -v github.com/gorilla/mux
   
3. Basic understanding of Graphql will be helpful.
   
4. It will be worth checking these sites for the structure of data returned by HSL HRT's open data framework.
   1. https://digitransit.fi/en/developers/apis/1-routing-api/
   2. https://api.digitransit.fi/graphiql/hsl

5. Google action supports mTLS. This means client and server communication can be secured using both server side and client side certificates and encryption keys. Details can be found here - https://cloud.google.com/dialogflow/docs/fulfillment-mtls.
   1. Let's Encrypt can be used to generate the server certificates to authenticate and authorize your webserver hosting this GO application - https://letsencrypt.org/
   2. Self-generated client certificate can also be generated for machines in development environment to run cURL commands during testing. This self generated certificate can be appended to ca-cert file that was generated for step-5-1 above for the Google servers.

## Deployment

1. Once all the GO packages are installed, build the application binary. For e.g. in a ubuntu shell: go build *.go

2. This creates a binary - ga-hsl-hrt. Run this application: ./ga-hsl-hrt
   
3. Check logfile for deployment status: For e.g. in a ubuntu shell: tail -f ./ga-hsl-hrt.log

## DialogFlow specifics
1. Application implements two intents with following names:
   1. Destination-Only: This intent is targetted for queries involving destination only. E.g:
      1. When is the next bus to Sello?
      2. Next route to Tapiola
   
   2. Bus-Destination: This intent is targetted for queries involving both bus and a destination. E.g:
      1. When is the next 215 to Sello?
      2. When is the next 321 to Helsinki?

2. Intent identifiers are defined in types.go

## Authors

* **Anand Radhakrishnan** - *Initial work* - [anand-p-r](https://github.com/anand-p-r)