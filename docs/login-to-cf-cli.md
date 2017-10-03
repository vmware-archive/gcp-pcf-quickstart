# Login to Cloud Foundry (CLI)

After you have [deployed PCF](./quick-deployment.md) you can access Cloud Foundry!

## Retrieve Credentials

From your deployment machine in the `gcp-pcf-quickstart` folder run the
following command:

```bash
bin/omg-cli remote --env-dir=env/omg "get-credential --app-name=cf --credential=.uaa.admin_credentials"
```

This will return an `Identity` and `Password`

## Target Cloud Foundry

```bash
cf login --skip-ssl-validation -a https://api.sys.$(util/terraform_output.sh dns_suffix)
```

Use the `Identity` from the previous setp as the `Email` and the
`Password` from the pervious step.

## Deploying Your First App

1. Clone the sample app:
   ```bash
   git clone https://github.com/cloudfoundry-samples/cf-sample-app-spring.git
   cd cf-sample-app-spring
   ```
1. Deploy the app:
   ```bash
   cf push
   ```

The last command will output the URL you can use to access the app.
