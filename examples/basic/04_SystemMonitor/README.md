# System Monitoring as a Service

Expose several system properties and logs.

## How to run it

```
$ kapow server SystemMonitor
```


## How to consume it

<details><summary>List files and directories</summary>

You can see if a file exists or the contents of a directory adding the path after _/file/_.

```
$ curl -s http://localhost:8080/file/var/log/cups
total 60
drwxr-xr-x  2 root root  4096 Dec  2 06:54 .
drwxr-xr-x 15 root root  4096 Jan 27 07:34 ..
-rw-r--r--  1 root cups     0 Dec  2 06:54 access_log
-rw-r--r--  1 root cups   125 Dec  1 19:12 access_log.1
-rw-r--r--  1 root cups   254 Nov 23 12:17 access_log.2
-rw-r--r--  1 root cups   125 Sep  8 14:41 access_log.3
-rw-r--r--  1 root cups   634 May  1  2019 access_log.4
-rw-r--r--  1 root cups     0 Sep  9  2018 error_log
-rw-r--r--  1 root cups 17312 Sep  3  2018 error_log.1
-rw-r--r--  1 root cups     0 Dec  2 06:54 page_log
-rw-r--r--  1 root cups   128 Dec  1 19:12 page_log.1
-rw-r--r--  1 root cups   188 Nov 23 12:17 page_log.2
-rw-r--r--  1 root cups   108 Sep  8 14:41 page_log.3
-rw-r--r--  1 root cups   465 May  1  2019 page_log.4
```
</details>

<details><summary>List processes</summary>

```
$ curl -s http://localhost:8080/process
nil        46717  0.0  0.0 111224  8196 pts/2    Sl   16:48   0:00 kapow server SystemMonitor
root       47405  0.0  0.0      0     0 ?        I    16:50   0:00 [kworker/3:1-mm_percpu_wq]
root       47406  0.0  0.0      0     0 ?        I    16:50   0:00 [kworker/0:1]
root       47819  0.0  0.0      0     0 ?        I    16:52   0:00 [kworker/7:2-mm_percpu_wq]
root       47823  0.6  0.0      0     0 ?        I<   16:52   0:01 [kworker/u17:0-hci0]
nil        48345  0.0  0.0   7124  2804 pts/2    S    16:56   0:00 /bin/sh -c ps -aux | kapow set /response/body
nil        48346  0.0  0.0   9392  3324 pts/2    R    16:56   0:00 ps -aux
nil        48347  0.0  0.0 109304  7080 pts/2    Sl   16:56   0:00 kapow set /response/body

...

```
</details>

<details><summary>CPU properties</summary>

```
$ curl -s http://localhost:8080/cpu | head -n 30
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 158
model name	: Intel(R) Xeon(R) CPU E3-1505M v6 @ 3.00GHz
stepping	: 9
microcode	: 0xca
cpu MHz		: 803.845
cache size	: 8192 KB
physical id	: 0
siblings	: 8
core id		: 0
cpu cores	: 4
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 22
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb invpcid_single pti ssbd ibrs ibpb stibp tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx rdseed adx smap clflushopt intel_pt xsaveopt xsavec xgetbv1 xsaves dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp md_clear flush_l1d
bugs		: cpu_meltdown spectre_v1 spectre_v2 spec_store_bypass l1tf mds swapgs taa itlb_multihit
bogomips	: 6002.00
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

```
</details>

<details><summary>Memory usage</summary>

```
$ curl -s http://localhost:8080/memory
              total        used        free      shared  buff/cache   available
Mem:          31967        3169       22636         734        6161       27608
Swap:             0           0           0

```
</details>

<details><summary>Disk usage</summary>

```
$ curl -s http://localhost:8080/disk/usage
Filesystem          Size  Used Avail Use% Mounted on
dev                  16G     0   16G   0% /dev
run                  16G  1.7M   16G   1% /run
/dev/nvme0n1p2      468G  419G   26G  95% /
tmpfs                16G  225M   16G   2% /dev/shm
tmpfs                16G     0   16G   0% /sys/fs/cgroup
tmpfs                16G  2.2M   16G   1% /tmp
/dev/nvme0n1p1      549M   62M  488M  12% /boot
/home/nil/.Private  468G  419G   26G  95% /home/nil
tmpfs               3.2G   12K  3.2G   1% /run/user/1000

```
</details>

<details><summary>Mount points</summary>

```
$ curl -s http://localhost:8080/disk/mounts
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
sys on /sys type sysfs (rw,nosuid,nodev,noexec,relatime)
dev on /dev type devtmpfs (rw,nosuid,relatime,size=16357376k,nr_inodes=4089344,mode=755)
run on /run type tmpfs (rw,nosuid,nodev,relatime,mode=755)
efivarfs on /sys/firmware/efi/efivars type efivarfs (rw,nosuid,nodev,noexec,relatime)
/dev/nvme0n1p2 on / type ext4 (rw,relatime)
securityfs on /sys/kernel/security type securityfs (rw,nosuid,nodev,noexec,relatime)
tmpfs on /dev/shm type tmpfs (rw,nosuid,nodev)
devpts on /dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=000)
tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755)
cgroup2 on /sys/fs/cgroup/unified type cgroup2 (rw,nosuid,nodev,noexec,relatime,nsdelegate)
cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,name=systemd)
pstore on /sys/fs/pstore type pstore (rw,nosuid,nodev,noexec,relatime)
none on /sys/fs/bpf type bpf (rw,nosuid,nodev,noexec,relatime,mode=700)
cgroup on /sys/fs/cgroup/perf_event type cgroup (rw,nosuid,nodev,noexec,relatime,perf_event)
cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_cls,net_prio)
cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpu,cpuacct)
cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
cgroup on /sys/fs/cgroup/rdma type cgroup (rw,nosuid,nodev,noexec,relatime,rdma)
cgroup on /sys/fs/cgroup/blkio type cgroup (rw,nosuid,nodev,noexec,relatime,blkio)
cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
cgroup on /sys/fs/cgroup/devices type cgroup (rw,nosuid,nodev,noexec,relatime,devices)
cgroup on /sys/fs/cgroup/hugetlb type cgroup (rw,nosuid,nodev,noexec,relatime,hugetlb)
cgroup on /sys/fs/cgroup/pids type cgroup (rw,nosuid,nodev,noexec,relatime,pids)
cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)

```
</details>

<details><summary>Open sockets</summary>

```
$ curl -s http://localhost:8080/socket
Netid   State    Recv-Q   Send-Q     Local Address:Port      Peer Address:Port  Process
udp     UNCONN   0        0           10.23.60.221:36922          0.0.0.0:*      users:(("chromium",pid=2005,fd=85))
udp     UNCONN   0        0            224.0.0.251:5353           0.0.0.0:*      users:(("chromium",pid=2031,fd=36))
udp     UNCONN   0        0            224.0.0.251:5353           0.0.0.0:*      users:(("opera",pid=1376,fd=139))
udp     UNCONN   0        0            224.0.0.251:5353           0.0.0.0:*      users:(("chromium",pid=2005,fd=133))
udp     UNCONN   0        0            224.0.0.251:5353           0.0.0.0:*      users:(("chromium",pid=2031,fd=172))
udp     UNCONN   0        0           10.23.60.221:51258          0.0.0.0:*      users:(("opera",pid=1376,fd=180))
tcp     LISTEN   0        4096           127.0.0.1:8081           0.0.0.0:*      users:(("kapow",pid=46717,fd=3))
tcp     LISTEN   0        4096           127.0.0.1:8082           0.0.0.0:*      users:(("kapow",pid=46717,fd=6))
tcp     LISTEN   0        5              127.0.0.1:631            0.0.0.0:*
tcp     LISTEN   0        4096                   *:8080                 *:*      users:(("kapow",pid=46717,fd=5))
tcp     LISTEN   0        5                  [::1]:631               [::]:*

```
</details>

<details><summary>Show kernel messages</summary>

```
$ curl -s http://localhost:8080/kernel/messages
[    0.319770] DMAR: ANDD device: 2 name: \_SB.PCI0.I2C1
[    0.319772] DMAR-IR: IOAPIC id 2 under DRHD base  0xfed91000 IOMMU 1
[    0.319772] DMAR-IR: HPET id 0 under DRHD base 0xfed91000
[    0.319773] DMAR-IR: Queued invalidation will be enabled to support x2apic and Intr-remapping.
[    0.321288] DMAR-IR: Enabled IRQ remapping in x2apic mode
[    0.321289] x2apic enabled
[    0.321302] Switched APIC routing to cluster x2apic.
[    0.325379] ..TIMER: vector=0x30 apic1=0 pin1=2 apic2=-1 pin2=-1
[    0.339725] clocksource: tsc-early: mask: 0xffffffffffffffff max_cycles: 0x2b3e459bf4c, max_idle_ns: 440795289890 ns
[    0.339730] Calibrating delay loop (skipped), value calculated using timer frequency.. 6002.00 BogoMIPS (lpj=10000000)
[    0.339732] pid_max: default: 32768 minimum: 301
[    0.342052] LSM: Security Framework initializing
[    0.342056] Yama: becoming mindful.
[    0.342119] Mount-cache hash table entries: 65536 (order: 7, 524288 bytes, linear)
[    0.342165] Mountpoint-cache hash table entries: 65536 (order: 7, 524288 bytes, linear)
[    0.342178] *** VALIDATE tmpfs ***
[    0.342325] *** VALIDATE proc ***
[    0.342373] *** VALIDATE cgroup1 ***
[    0.342374] *** VALIDATE cgroup2 ***
[    0.342425] mce: CPU0: Thermal monitoring enabled (TM1)

```
</details>

<details><summary>Show systemd journal</summary>

```
$ curl -s http://localhost:8080/systemd/journal | head -n 10
-- Logs begin at Mon 2019-07-29 07:35:17 CEST, end at Mon 2020-01-27 17:05:10 CET. --
Jan 27 07:34:37 xan kernel: hid-generic 0003:1D50:6122.0008: input,hidraw7: USB HID v1.10 Mouse [Ultimate Gadget Laboratories Ultimate Hacking Keyboard] on usb-0000:0c:00.0-1.7/input4
Jan 27 07:34:37 xan mtp-probe[1053]: checking bus 3, device 5: "/sys/devices/pci0000:00/0000:00:1c.4/0000:04:00.0/0000:05:01.0/0000:07:00.0/0000:08:04.0/0000:0a:00.0/0000:0b:01.0/0000:0c:00.0/usb3/3-1/3-1.6"
Jan 27 07:34:37 xan mtp-probe[1055]: checking bus 3, device 3: "/sys/devices/pci0000:00/0000:00:1c.4/0000:04:00.0/0000:05:01.0/0000:07:00.0/0000:08:04.0/0000:0a:00.0/0000:0b:01.0/0000:0c:00.0/usb3/3-1/3-1.1"
Jan 27 07:34:37 xan mtp-probe[1054]: checking bus 3, device 6: "/sys/devices/pci0000:00/0000:00:1c.4/0000:04:00.0/0000:05:01.0/0000:07:00.0/0000:08:04.0/0000:0a:00.0/0000:0b:01.0/0000:0c:00.0/usb3/3-1/3-1.7"
Jan 27 07:34:37 xan mtp-probe[1053]: bus: 3, device: 5 was not an MTP device
Jan 27 07:34:37 xan mtp-probe[1054]: bus: 3, device: 6 was not an MTP device
Jan 27 07:34:37 xan mtp-probe[1056]: checking bus 3, device 4: "/sys/devices/pci0000:00/0000:00:1c.4/0000:04:00.0/0000:05:01.0/0000:07:00.0/0000:08:04.0/0000:0a:00.0/0000:0b:01.0/0000:0c:00.0/usb3/3-1/3-1.5"
Jan 27 07:34:37 xan mtp-probe[1056]: bus: 3, device: 4 was not an MTP device
Jan 27 07:34:37 xan mtp-probe[1055]: bus: 3, device: 3 was not an MTP device

```
</details>
