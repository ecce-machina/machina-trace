This is the idea, have ways to inquire about a compute side workload signature with the following data points

               +-------------------+
               |  machina-agent    |
               +-------------------+
                        |
        +---------------+---------------+
        |               |               |
     procfs         filesystem        eBPF
        |               |               |
   +---------+   +------+------+   +---------+
   | vmstat  |   | Lustre | Ceph |   | syscalls|
   | meminfo |   | BeeGFS | NFS  |   | VFS     |
   | disks   |   | GPFS   | ...  |   | block   |
   +---------+   +------+------+   +---------+
                        |
                  unified Snapshot
