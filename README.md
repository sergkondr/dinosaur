# dinosaur

Simple DNS server wannabe. The main goal of this project was to practice working with raw network data and enjoy the process. 
I chose DNS because it is a simple and clear protocol with an easy-to-understand format.

It works, and I am happy with the result, that means the goal is achieved.

Maybe one day, I will develop this project further.

### Usage

```shell
➜ go build cmd/* && ./dinosaur
Dinosaur is listening on 127.0.0.1:9053
DNS query ID: 37627 Questions: [test example com]

found following addresses for  [test example com] : [[127 0 0 2] [127 0 0 3]]
creating response for  [test example com]  with ID 37627  and addresses [[127 0 0 2] [127 0 0 3]]
```
```shell
➜ dig @127.0.0.1 -p 9053 test.example.com

; <<>> DiG 9.10.6 <<>> @127.0.0.1 -p 9053 test.example.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 37627
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;test.example.com.              IN      A

;; ANSWER SECTION:
test.example.com.       20      IN      A       127.0.0.2
test.example.com.       20      IN      A       127.0.0.3

;; Query time: 0 msec
;; SERVER: 127.0.0.1#9053(127.0.0.1)
;; WHEN: Mon Jan 06 13:33:43 CET 2025
;; MSG SIZE  rcvd: 66
```
