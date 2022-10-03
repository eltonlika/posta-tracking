# posta-tracking
TUI client for Albanian Post tracking service

Build with ```go build``` or ```go install``` to $GOPATH/bin

```
Usage: ./posta-tracking <tracking number>
```


Example output:
```
» posta-tracking RB560254054SG
#  Kodi           Data                 Ngjarja                                                                       Zyra                                 Destinacioni
1  RB560254054SG  2018-06-12 02:39 AM  Objekti u fut ne thes / Insert item into bag (Otb)                            Singapor
2  RB560254054SG  2018-06-21 16:11 PM  U pranua objekt nga jashte / Receive item from abroad (EDI-received)          Posta e Jashtme (mberritje), TIRANA
3  RB560254054SG  2018-06-21 16:42 PM  Objekti u fut ne thes / Insert item into bag (Otb)                            Posta e Jashtme (mberritje), TIRANA  ZP Tirana 5/1, TIRANA
4  RB560254054SG  2018-06-21 17:30 PM  Objekti u dergua ne zyren destinacion / Send item to domestic location (Inb)  Posta e Jashtme (mberritje), TIRANA  Dispeceria Tirane, TIRANA
5  RB560254054SG  2018-06-21 17:35 PM  Mberritje objekti ne destinacion / Receive item at location (Otb)             Dispeceria Tirane, TIRANA
6  RB560254054SG  2018-06-21 17:36 PM  Skanuar per transport / Scan to transport                                     Dispeceria Tirane, TIRANA            ZP Tirana 5/1, TIRANA
7  RB560254054SG  2018-06-24 09:04 AM  Objekti u nxor nga thesi per perpunim / Receive item at location (Inb)        ZP Tirana 5/1, TIRANA
8  RB560254054SG  2018-06-29 20:06 PM  Objekti u dorezua / Deliver item (Inb)                                        ZP Tirana 5/1, TIRANA
```
