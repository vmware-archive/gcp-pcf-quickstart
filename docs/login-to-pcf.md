# Login to Pivotal Cloud Foundry

After you have [deployed PCF](./quick-deployment.md) it's time to use it!

## Access Pivotal Operations Manager (Ops Man)

Run the following command from the `gcp-pcf-quickstart` folder on your deployment machine.
```bash
util/env_info.sh opsman
```

The command will output the URL, username, password to access the Ops Manager web interface.

## <a name="cfapi"></a>Access the Cloud Foundry API

### 1. Retrieve Credentials

Run the following command from the `gcp-pcf-quickstart` folder on your deployment machine:
```bash
util/env_info.sh cf
```

This will return an `identity` and `password`

### 2. Target Cloud Foundry

```bash
cf login --skip-ssl-validation -a https://api.sys.$(ENV_DIR=env/pcf util/terraform_output.sh dns_suffix)
```

Use the `identity` from the previous step as the `Email` and the
`password` from the previous step. Select `system` for space when prompted.

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

## Access the Jumpbox

Run the following command from the `gcp-pcf-quickstart` folder on your deployment machine.

  ```bash
  util/ssh.sh
  ```

This command will open an SSH session with the jumpbox deployed to your PCF network.

## What's Next?
- [Deploy an app using Google Cloud Storage and Vision APIs](./deploy-awwvision.md)
- [Delete Deployment](./deleting-deployment.md)
