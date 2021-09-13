server sends:

0: Connection Confirmed
1: State Packet 4B gametick, 4B x 4B y
2: State Packet 4B gametick, client amount, for each client: (2B sid 4B x 4B y)

client sends:
0: Connect request
1: State Packet 4B gametick 4B x 4B y
