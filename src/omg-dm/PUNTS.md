# Punted Work

Debt that needs to be resolved but not at this moment.

## NAT per zone

NAT is deployed to a single zone. This will cause cross-zone access for NAT. 

The tricky thing here is: how does PCF/tag instances?

## GCE Enforcer

Currently giving the enforcer an IAM binding on the project. This is Googler specific.

## Global Resources

Move IAM roles (ops_manager/ops_manager.jinja) to a global roles or base.

## DNS

Deployment Manager can't create individual DNS records, only managed zones
