# Punted Work

Debt that needs to be resolved but not at this moment.

## NAT per zone

A sinlge NAT instance group is being deployed. All traffic will cross zones to this group. This should be replaced by a NAT instance group per zone.

The tricky thing here is: how does PCF/tag instances?