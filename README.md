# Astra
On Demand Log Monitoring Service

[Project Brief](https://docs.google.com/document/d/1a44aCd8sJOLXWNT2vTerRGuNAYSHsTLxdFJWxRCl2pA/edit#)

[Project Tech Design](https://docs.google.com/document/d/1vjoh4NN57tht9tUDZSYe5K_T2J2nhUxp940qWcUisOw/edit#heading=h.2vb45ps60q20)

[Github Project Board](https://github.com/users/karthikbic1/projects/1)


# Development

Run the following command to get started..

```

# Install Dependencies for the first time.
$ bash scripts/dev.sh init

# Start the Astra App - defaults on port 3000
$ bash scripts/dev.sh up

# Start the Astra App on a different port
$ bash scripts/dev.sh up 3001

# To run unit tests
$ bash scripts/dev.sh test

```

# Examples

### Request logs with default parameters (default lines: 10)

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=syslog" -i
HTTP/1.1 200 OK
Date: Tue, 04 Oct 2022 06:57:37 GMT
Content-Length: 1471
Content-Type: text/plain; charset=utf-8

{"Logs":["Oct  3 23:57:36 karthik-Lenovo-IdeaPad-S145-15AST systemd[1]: tmp-snap.rootfs_xOyYnH.mount: Succeeded.","Oct  3 23:57:36 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: tmp-snap.rootfs_xOyYnH.mount: Succeeded.","Oct  3 23:57:36 karthik-Lenovo-IdeaPad-S145-15AST kernel: [  470.057703] audit: type=1400 audit(1664866656.776:52): apparmor=\"DENIED\" operation=\"capable\" profile=\"/snap/snapd/17029/usr/lib/snapd/snap-confine\" pid=4000 comm=\"snap-confine\" capability=4  capname=\"fsetid\"","Oct  3 23:57:36 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started snap.curl.curl.04fbe18e-5957-4a1f-9ee4-b34b9313b262.scope.","Oct  3 23:57:34 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started VTE child process 3991 launched by gnome-terminal-server process 3814.","Oct  3 23:57:19 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: snap.go.go.715d0e3c-663f-4d91-82ee-dc54873cf5b4.scope: Succeeded.","Oct  3 23:57:19 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started snap.go.go.715d0e3c-663f-4d91-82ee-dc54873cf5b4.scope.","Oct  3 23:57:13 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: snap.go.go.d4cc35bc-aac0-41a7-a390-a40684504203.scope: Succeeded.","Oct  3 23:57:13 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started VTE child process 3907 launched by gnome-terminal-server process 3814.","Oct  3 23:57:11 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started snap.go.go.d4cc35bc-aac0-41a7-a390-a40684504203.scope."],"ErrorMsg":""}
```

### Request logs with num_lines param

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=syslog&num_lines=2" -i
HTTP/1.1 200 OK
Date: Tue, 04 Oct 2022 06:58:43 GMT
Content-Length: 296
Content-Type: text/plain; charset=utf-8

{"Logs":["Oct  3 23:58:43 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: Started snap.curl.curl.9ec4d6cb-502a-4e8d-9907-4707f6026241.scope.","Oct  3 23:57:37 karthik-Lenovo-IdeaPad-S145-15AST systemd[1342]: snap.curl.curl.04fbe18e-5957-4a1f-9ee4-b34b9313b262.scope: Succeeded."],"ErrorMsg":""}

```

### Request logs with filter param

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=syslog&num_lines=2000000&filter=DENIED" -i
HTTP/1.1 200 OK
Date: Tue, 04 Oct 2022 07:04:24 GMT
Content-Length: 287
Content-Type: text/plain; charset=utf-8

{"Logs":["Oct  4 00:00:09 karthik-Lenovo-IdeaPad-S145-15AST kernel: [  622.543007] audit: type=1400 audit(1664866809.258:53): apparmor=\"DENIED\" operation=\"capable\" profile=\"/usr/sbin/cups-browsed\" pid=4187 comm=\"cups-browsed\" capability=23  capname=\"sys_nice\""],"ErrorMsg":""}

```

### Request logs with seconddary server.

1. Start two Astra servers in two seperate terminals

```
bash script/dev up #terminal 1
bash script/dev up 3001 #terminal 2
```

2. Create a file `/var/log/testlog-example-seconday` and ensure that `/var/log/testlog-example` doesnt exist to mimic primary server not having the log file.

3. Make a curl call

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=testlog-example&secondary_server=http://localhost:3001/" -i
HTTP/1.1 200 OK
Date: Tue, 04 Oct 2022 07:05:04 GMT
Content-Length: 192
Content-Type: text/plain; charset=utf-8

{"Logs":["this is 9 line.","this is 8 line ","this is 7 line","this is 6 line.","this is 5 line 今日は","this is 4 line","this is 3 line.","this is 2 line","this is 1 line"],"ErrorMsg":""}


```
Also dont provide the secondary server 

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=testlog-example" -i
HTTP/1.1 404 Not Found
Date: Tue, 04 Oct 2022 07:06:27 GMT
Content-Length: 152
Content-Type: text/plain; charset=utf-8

{"Logs":null,"ErrorMsg":"Primary Server:open /var/log/testlog-example: no such file or directory, Secondary Server:Not fetching from secondary server"}

```