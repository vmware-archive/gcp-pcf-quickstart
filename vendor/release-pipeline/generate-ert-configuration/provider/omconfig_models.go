package provider

type Property struct {
	Value interface{} `json:"value,omitempty"`
}

type ERTProperties struct {
	SysDomain                    *Property `json:".cloud_controller.system_domain,omitempty"`
	AppsDomain                   *Property `json:".cloud_controller.apps_domain,omitempty"`
	HAProxySkipSSL               *Property `json:".ha_proxy.skip_cert_verify,omitempty"`
	MySQLRecipientEmail          *Property `json:".mysql_monitor.recipient_email,omitempty"`
	SAMLCert                     *Property `json:".uaa.service_provider_key_credentials,omitempty"`
	NetworkingPOE                *Property `json:".properties.networking_point_of_entry,omitempty"`
	ExternalSSLCert              *Property `json:".properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate,omitempty"`
	HAProxySSLCert               *Property `json:".properties.networking_point_of_entry.haproxy.ssl_rsa_certificate,omitempty"`
	SecurityCheckBox             *Property `json:".properties.security_acknowledgement,omitempty"`
	SystemDB                     *Property `json:".properties.system_database,omitempty"`
	ExtDBPort                    *Property `json:".properties.system_database.external.port,omitempty"`
	ExtDBHost                    *Property `json:".properties.system_database.external.host,omitempty"`
	AppUsageDBUsername           *Property `json:".properties.system_database.external.app_usage_service_username,omitempty"`
	AppUsageDBPassword           *Property `json:".properties.system_database.external.app_usage_service_password,omitempty"`
	AutoscaleDBUsername          *Property `json:".properties.system_database.external.autoscale_username,omitempty"`
	AutoscaleDBPassword          *Property `json:".properties.system_database.external.autoscale_password,omitempty"`
	CCDBUsername                 *Property `json:".properties.system_database.external.ccdb_username,omitempty"`
	CCDBPassword                 *Property `json:".properties.system_database.external.ccdb_password,omitempty"`
	DiegoDBUsername              *Property `json:".properties.system_database.external.diego_username,omitempty"`
	DiegoDBPassword              *Property `json:".properties.system_database.external.diego_password,omitempty"`
	NetworkPolicyDBUsername      *Property `json:".properties.system_database.external.networkpolicyserver_username,omitempty"`
	NetworkPolicyDBPassword      *Property `json:".properties.system_database.external.networkpolicyserver_password,omitempty"`
	NFSVolumeDBUsername          *Property `json:".properties.system_database.external.nfsvolume_username,omitempty"`
	NFSVolumeDBPassword          *Property `json:".properties.system_database.external.nfsvolume_password,omitempty"`
	NotificationsDBUsername      *Property `json:".properties.system_database.external.notifications_username,omitempty"`
	NotificationsDBPassword      *Property `json:".properties.system_database.external.notifications_password,omitempty"`
	AccountDBUsername            *Property `json:".properties.system_database.external.account_username,omitempty"`
	AccountDBPassword            *Property `json:".properties.system_database.external.account_password,omitempty"`
	RoutingDBUsername            *Property `json:".properties.system_database.external.routing_username,omitempty"`
	RoutingDBPassword            *Property `json:".properties.system_database.external.routing_password,omitempty"`
	UAADBUsername                *Property `json:".properties.system_database.external.uaa_username,omitempty"`
	UAADBPassword                *Property `json:".properties.system_database.external.uaa_password,omitempty"`
	SystemBlobstore              *Property `json:".properties.system_blobstore,omitempty"`
	ExternalPackagesBucket       *Property `json:".properties.system_blobstore.external.packages_bucket,omitempty"`
	ExternalDropletsBucket       *Property `json:".properties.system_blobstore.external.droplets_bucket,omitempty"`
	ExternalResourcesBucket      *Property `json:".properties.system_blobstore.external.resources_bucket,omitempty"`
	ExternalBuildpacksBucket     *Property `json:".properties.system_blobstore.external.buildpacks_bucket,omitempty"`
	StorageIAMEndpoint           *Property `json:".properties.system_blobstore.external.endpoint,omitempty"`
	StorageIAMAccessKey          *Property `json:".properties.system_blobstore.external.access_key,omitempty"`
	StorageIAMSSecretAccessKey   *Property `json:".properties.system_blobstore.external.secret_key,omitempty"`
	StorageIAMRegion             *Property `json:".properties.system_blobstore.external.region,omitempty"`
	GCSPackagesBucket            *Property `json:".properties.system_blobstore.external_gcs.packages_bucket,omitempty"`
	GCSDropletsBucket            *Property `json:".properties.system_blobstore.external_gcs.droplets_bucket,omitempty"`
	GCSResourcesBucket           *Property `json:".properties.system_blobstore.external_gcs.resources_bucket,omitempty"`
	GCSBuildpacksBucket          *Property `json:".properties.system_blobstore.external_gcs.buildpacks_bucket,omitempty"`
	GCSAccessKey                 *Property `json:".properties.system_blobstore.external_gcs.access_key,omitempty"`
	GCSSecretKey                 *Property `json:".properties.system_blobstore.external_gcs.secret_key,omitempty"`
	TCPRouting                   *Property `json:".properties.tcp_routing,omitempty"`
	TCPReservablePorts           *Property `json:".properties.tcp_routing.enable.reservable_ports,omitempty"`
	SMTPFrom                     *Property `json:".properties.smtp_from,omitempty"`
	SMTPAddress                  *Property `json:".properties.smtp_address,omitempty"`
	SMTPPort                     *Property `json:".properties.smtp_port,omitempty"`
	SMTPCredentials              *Property `json:".properties.smtp_credentials,omitempty"`
	SMTPEnableStartTLSAuto       *Property `json:".properties.smtp_enable_starttls_auto,omitempty"`
	SMTPAuthMechanism            *Property `json:".properties.smtp_auth_mechanism,omitempty"`
	LoggerEndpointPort           *Property `json:".properties.logger_endpoint_port,omitempty"`
	CFStorageAccountName         *Property `json:".properties.system_blobstore.external_azure.account_name,omitempty"`
	CFStorageAccountAccessKey    *Property `json:".properties.system_blobstore.external_azure.access_key,omitempty"`
	CFBuildpacksStorageContainer *Property `json:".properties.system_blobstore.external_azure.buildpacks_container,omitempty"`
	CFDropletsStorageContainer   *Property `json:".properties.system_blobstore.external_azure.droplets_container,omitempty"`
	CFPackagesStorageContainer   *Property `json:".properties.system_blobstore.external_azure.packages_container,omitempty"`
	CFResourcesStorageContainer  *Property `json:".properties.system_blobstore.external_azure.resources_container,omitempty"`
	ContainerNetworking          *Property `json:".properties.container_networking,omitempty"`
}

type Network struct {
	Name string `json:"name,omitempty"`
}

type AZ struct {
	Name string `json:"name,omitempty"`
}

type ERTNetworks struct {
	SingletonAZ AZ      `json:"singleton_availability_zone,omitempty"`
	OtherAZs    []AZ    `json:"other_availability_zones,omitempty"`
	Network     Network `json:"network,omitempty"`
}

type JobResourceConfig struct {
	Instances         int      `json:"instances,omitempty"`
	InternetConnected bool     `json:"internet_connected,omitempty"`
	ELBNames          []string `json:"elb_names,omitempty"`
}

type ERTResources struct {
	TCPRouter     *JobResourceConfig `json:"tcp_router,omitempty"`
	HAProxy       *JobResourceConfig `json:"ha_proxy,omitempty"`
	MySQL         *JobResourceConfig `json:"mysql,omitempty"`
	Router        *JobResourceConfig `json:"router,omitempty"`
	ConsulServer  *JobResourceConfig `json:"consul_server,omitempty"`
	EtcdTLSServer *JobResourceConfig `json:"etcd_tls_server,omitempty"`
	EtcdServer    *JobResourceConfig `json:"etcd_server,omitempty"`
	DiegoBrain    *JobResourceConfig `json:"diego_brain,omitempty"`
	DiegoCell     *JobResourceConfig `json:"diego_cell,omitempty"`
	DiegoDB       *JobResourceConfig `json:"diego_database,omitempty"`
	MySQLProxy    *JobResourceConfig `json:"mysql_proxy,omitempty"`
}

type ERTConfiguration struct {
	OpsManagerDomain  string        `json:"ops_manager_domain"`
	ProductProperties ERTProperties `json:"properties,omitempty"`
	ProductNetwork    ERTNetworks   `json:"network,omitempty"`
	ProductResources  ERTResources  `json:"resources,omitempty"`
}

type FullConfiguration struct {
	ERTConfig ERTConfiguration `json:"product,omitempty"`
}
