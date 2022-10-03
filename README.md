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
Date: Tue, 20 Sep 2022 05:58:44 GMT
Content-Length: 1205
Content-Type: text/plain; charset=utf-8

{"Logs":"\n\nSep 19 22:58:44 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.curl.curl.ca9af74e-29a4-41ba-822a-2545730c6dab.scope.\nSep 19 22:58:27 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: snap.go.go.c496bcfa-8ae3-4380-8847-e30d83d1de96.scope: Succeeded.\nSep 19 22:58:26 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.go.go.c496bcfa-8ae3-4380-8847-e30d83d1de96.scope.\nSep 19 22:58:21 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: snap.go.go.e19334c3-92c6-4edd-9af6-d4bfe86bf66e.scope: Succeeded.\nSep 19 22:58:21 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.go.go.e19334c3-92c6-4edd-9af6-d4bfe86bf66e.scope.\nSep 19 22:58:17 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: snap.go.go.0dc13ce4-47ba-4281-8936-524a3dd0b013.scope: Succeeded.\nSep 19 22:58:17 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.go.go.0dc13ce4-47ba-4281-8936-524a3dd0b013.scope.\nSep 19 22:58:14 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: snap.go.go.69ff5d91-9e5e-4b3e-93d1-88d315ec97dd.scope: Succeeded.\nSep 19 22:58:14 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.go.go.69ff5d91-9e5e-4b3e-93d1-88d315ec97dd.scope.","ErrorMsg":""}
```

### Request logs with num_lines param

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=syslog&num_lines=2" -i
HTTP/1.1 200 OK
Date: Tue, 20 Sep 2022 05:45:47 GMT
Content-Length: 161
Content-Type: text/plain; charset=utf-8

{"Logs":"\n\nSep 19 22:45:47 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.curl.curl.6ad21df8-0318-4eed-810f-7ad3f73a9921.scope.","ErrorMsg":""}

```

### Request logs with filter param

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=syslog&num_line=100&filter=snap.go.go.0efe9287-5188-4f42-96c5-f915f38c451b.scope" -i
HTTP/1.1 200 OK
Date: Tue, 20 Sep 2022 05:50:35 GMT
Content-Length: 287
Content-Type: text/plain; charset=utf-8

{"Logs":"\nSep 19 22:48:43 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: snap.go.go.0efe9287-5188-4f42-96c5-f915f38c451b.scope: Succeeded.\nSep 19 22:48:43 karthik-Lenovo-IdeaPad-S145-15AST systemd[1475]: Started snap.go.go.0efe9287-5188-4f42-96c5-f915f38c451b.scope.","ErrorMsg":""}
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
Date: Tue, 20 Sep 2022 07:03:23 GMT
Content-Length: 186
Content-Type: text/plain; charset=utf-8

{"Logs":"\nthis is 9 line.\nthis is 8 line \nthis is 7 line\nthis is 6 line.\nthis is 5 line 今日は\nthis is 4 line\nthis is 3 line.\nthis is 2 line\nthis is 1 line\n","ErrorMsg":""}


```
Also dont provide the secondary server 

```
curl -X GET "http://localhost:3000/fetchlogs?file_name=testlog-example" -i
HTTP/1.1 404 Not Found
Date: Tue, 20 Sep 2022 19:32:39 GMT
Content-Length: 150
Content-Type: text/plain; charset=utf-8

{"Logs":"","ErrorMsg":"Primary Server:open /var/log/testlog-example: no such file or directory, Secondary Server:Not fetching from secondary server"}

```