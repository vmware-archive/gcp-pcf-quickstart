# Login to Pivotal Cloud Foundry

After you have [deployed PCF](./quick-deployment.md) your installation!

## Access Pivotal Operations Manager

Run the following command from the `gcp-pcf-quickstart` folder on your deployment machine.
```bash
printf "url: https://$(util/terraform_output.sh ops_manager_dns)\nusername: $(util/terraform_output.sh ops_manager_username)\npassword: $(util/terraform_output.sh ops_manager_password)\n)"
```

The command will output the URL, username, password to access the Ops Manager web interface.

## Access the Cloud Foundry API

### 1. Retrieve Credentials

Run the following command from the `gcp-pcf-quickstart` folder on your deployment machine:
```bash
bin/omg-cli remote --env-dir=env/omg "get-credential --app-name=cf --credential=.uaa.admin_credentials"
```

This will return an `Identity` and `Password`

### 2. Target Cloud Foundry

```bash
cf login --skip-ssl-validation -a https://api.sys.$(util/terraform_output.sh dns_suffix)
```

Use the `Identity` from the previous setp as the `Email` and the
`Password` from the pervious step.

### 3. Deploying Your First App

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
