# Balancers
This directory contains the standard load balancers that MSTK provides out of the box.

## ```OneToOne```
In some cases a load balancer is not needed in such case to keep the architecture of MSTK intact a OneToOne balancer is used. This allows
to keep the same architecture all the while keeping an option in the future changing it for another balancer by just running a different
one.

The 'OneToOne' balancer just forwards requests straight to a shard.
