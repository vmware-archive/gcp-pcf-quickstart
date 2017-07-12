# Done
- network: ops man, ert (networks.tf, cf_internal_network.tf)
- jumpbox (not included)
- OpsMan VM & iam (ops_manager.tf, iam.tf)
- NAT (not included)

# Untested
- routing (http, https, ssh)
- storage.tf
- tcp_router.tf

# TODO
- output/parameterize `omg-opsman`, use it for next-hop on NAT rules.

# Next Steps
- IAM is messed up. do it out of the script.
- isolation segment (isolation_segment/)
- cloud SQL for OpsMan/ERT (external_database/)
- service broker: IAM, Cloud SQL
- stackdriver-tools: IAM

