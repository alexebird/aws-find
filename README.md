aws-find
========

ec2
---

```
$ af ec2 -f foobar -c
PRIVATE_IP  NAME         STATE    TYPE      IMAGE         KEY
10.0.1.8    foobar-host  running  t2.micro  ami-12345678  default

==> connecting to foobar-host(10.0.1.8)
==> via command: ssh 10.0.1.8

Welcome to Ubuntu 16.04.3 LTS (GNU/Linux 4.4.0-1038-aws x86_64)
...
```
