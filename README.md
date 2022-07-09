# P2PTunnel
peer to peer based tunnel

## Terminaology
P2P tunnel has 2 running parts to connect your two different devices:
- Agent: running at your home to forward packets to your local services.
- Connector: running at your remote device to connect to agent through p2p network, and proxy client request to agent.

## Usage

```
NAME:
   p2ptunnel - p2p tunnel

USAGE:
   p2ptunnel [global options] command [command options] [arguments...]

AUTHOR:
    <dwebfan@gmail.com>

COMMANDS:
   add, a     add peer name and its ID
   agent      start p2p tunnel agent service
   connector  start p2p tunnel connector service
   init, i    user friendly name of agent or connection
   remove, a  remove peer name and its ID
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --conf value, -c value  config file path (default: "./conf/p2ptunnel.conf")
   --help, -h              show help
```

## Quick start
1. Init agent config file at your home device. Default configuration file is ./conf/p2ptunnel.yml.

```
[agent-node] $ p2ptunnel init agent
Initialized new config at ./conf/p2ptunnel.yml
Please remember your ID: QmW1RE2V9aprXPKGXn8SqBzj34egVLaGeaVZdmZeJtthK6
```

You can provide your own config file name and path by `p2ptunnel -c <your path> init agent`

2. Init connector config file at your another device.

```
[connector-node] $ p2ptunnel init connection
Initialized new config at ./conf/p2ptunnel.yml
Please remember your ID: QmdiDf3DRWhkUVz5hwhC8ax7PHW1EmzdSDRE2JUx8TDucy
```

3. Add connector ID at your home device.
```
[agent-node] $ ./p2ptunnel add connector QmdiDf3DRWhkUVz5hwhC8ax7PHW1EmzdSDRE2JUx8TDucy
connector - QmdiDf3DRWhkUVz5hwhC8ax7PHW1EmzdSDRE2JUx8TDucy has been saved in config file: ./conf/p2ptunnel.yml
```
Now, config file should have below contents
```
[agent-node] $ cat ./conf/p2ptunnel.yml
name: agent
id: QmW1RE2V9aprXPKGXn8SqBzj34egVLaGeaVZdmZeJtthK6
private_key: !!binary |
  CAASqQkwggSlAgEAAoIBAQC8ANUjGBUco38nQJ3rnP5FhQFQtA4XCFSf++4bmwjOL3HMQE
  adyPWto2vxP5nafIOdBLeisxSZkneffcr95fJFYrENcilESWq+2JlwjIpQINA2KnLEJpK0
  kq7rIzqnRChfm3Ph5HuL1FycpISLuQz9I5lPuR3E1b9fPQCwhjOoQ2NE+87z5a2ZNnbwPW
  vSseroSiSRvN5/fD0BZyt8m8MgiBSFP9XB9oJfXGoC8r8u2z4g2GKsnrcJZVYnyCW3EdSB
  1rw++K0nQTFWxT6u0CPuO8BKvJmzW1DVD5p+ENCakgINhPMP7ie6sTmtJ+wvj1XuyIgR5z
  yLFE2HykvC2rRjAgMBAAECggEAfiBU3SFn8HmjcXMBPzNoxrzvX/Qby8nz9Ayw9mYxJxpf
  TvcEKRDL+Xy8ivHvRsvoBCxJAmb/9/NyO1bKG0HsIi6Ot4WSN9TFP1nSvtYaaJ1K8jvSGN
  QD8g7COM++pr6sE1AYE0EUvX9vzkr0/UPdeeorDXgUh5wasksbrlGqUEQWGa1IfPcdpwHn
  V4rmKKxKYSmms+8kP3nLVTWEFORfJfsoDXx5UR4Z3/vT1mKK2bF1tq3YfWyf7HCmZ4+BsK
  GcSkKnicH4C90i3NvLmMgCcxEdf4vl14n9iv8AKU1K30tkkm10MvGfvOPx5K8N0822aqp8
  S+1B0ik05FbiMPJawQKBgQDszYQYcmo6oS1qEIaaqh/C+zrv95Urp96q+0aUyekVytY/6K
  awTow5DQbCliTFawqNsKvWx0Gwnsyng39BDwB1xb9666WfATLe1WNyc6mf7F8el/yhTBBm
  AByOK9P6PjLO7axrvR620sr9NScFnLjtbszRfigKOHcW9ZxV97qcgwKBgQDLPo5JNZNyMB
  1qYCb++QhzzwPOQvsUUF59cnBxd16O+VLvNhDe5URRT0YOcq1lqUM3ZHSuWK2kXUJqjTdg
  pBzSzth4yGm2Qh0edSmYi/g/kwy0g0s23YDFA7OI7AmiAcKEBKMLeztCT5JaHQ+6TLyZ9F
  d1lVzAGW5ChAU49SHCoQKBgQCT5QtiRVspQO7fNnEK+cOagFPf+a41tMNx4DvFw5EKpkNH
  aONqa7RSEVuufh17Gw5dTgEUxB+30oYY/RymIlt0MswTVkd7VkFSQM26dphzJCqILf5/Ms
  VvHxS2ipL60IvlBzXPmC8tmdtjZyX28FnjGHddQ8B4GanvMfixDGaFRwKBgQCdQA8ykWM9
  TADWVwKU7/UcNSVKpwRAWVZiVPKut57PnBQQxJIVAunyxxT7BLsoFufMqcqlQjNHImjKq+
  wWt6Mdb7CI8vbnbwu8jwXZ4yH1fj6sQ5EkKAkDZbO40nc5g4cOQdAsh/H3gj1Hv4h1qf2a
  WDR6409ZydNHX4Hy6aZvQQKBgQCqUsmqEoz80xGdBhr0H3138nhps4QoRpGEsFj6m8cODX
  C909bQdL+smdyNPpM16DaRdj/BogENWIEa2pGb6SFFm12/FgHuifFIJgDp2Ro9j1eAPDLj
  icPsoWmMgeikvSnFiv4+deCnH1ot0OVjBV6qa87pXyubFZEbE0+TNbAPsg==
peers:
  connector:
    id: QmdiDf3DRWhkUVz5hwhC8ax7PHW1EmzdSDRE2JUx8TDucy
```
4. Add agent ID at your connector node
```
[connector-node] $ ./p2ptunnel add agent QmW1RE2V9aprXPKGXn8SqBzj34egVLaGeaVZdmZeJtthK6
connector - QmW1RE2V9aprXPKGXn8SqBzj34egVLaGeaVZdmZeJtthK6 has been saved in config file: ./conf/p2ptunnel.yml
```
Now, config file should have below contents
```
[connector-node] $ cat ./conf/p2ptunnel.yml
name: connector
id: QmdiDf3DRWhkUVz5hwhC8ax7PHW1EmzdSDRE2JUx8TDucy
private_key: !!binary |
  CAASpgkwggSiAgEAAoIBAQCwMQna7S6RuXsyc7ToiDYy7DApHwCEkGtwluxnR4kHDTc7kK
  M5kA0Hk2kZa2vq8wkZ4HfsG3sVmggWiK+ajV+HHcYsnu+HGI6MAjYsnJjhotztNIjvZIzK
  t4lln1xBk90O97xFBtCZ6JFicVK7T4Z2uJ297hh2WlSkV1nsHDHPecTDl1m/0CSF9/1H+/
  lvUJkS4v4VPM2+ETJRxmTlKiC25u7fsudKl+zmHTXTkPVvpSXQJljD+30/RyrgIslw4yPs
  11fr5ZZ9pXLNdKFZqVPh1yxLLswCPX15KhHpswp8ugQ81h8VcRX2RWSxm3jyio9eeqqLsl
  c0lFRyXQxbphlvAgMBAAECggEAAtz5JQafsByMhPheYzz7bH5sFe78Cityo4TAWLlP3752
  PFCQZnoRzCK4HYKiYVILvtDoAf08VdCH+x3DhMZxW/e+5bC7gb2Da4EJslXlIh2Ma4pkA0
  fmBdFPuUgKrsIhYIHkHFcNAsNFwYzH5GVZcQp0/cYlvZ6gK3+D5ZNbt9xiy3/pISAE1kwM
  elRa+YaWuOv61ABIxUgobLIcjAquq5r6j9USGs1DcgsC6GlW8ZlNLn8s7auwR/DUPfCf5J
  wwRFLpBGkIK2ZWfzLLsaROmTPgTOXPIBpadlL3RMrB3jsHq9U4GLmSl3ml5d8QXjiS3Ii2
  fhy3wOGyEAuBcgXEQQKBgQDkTWTDwCCNQ8Mdi3Yg9ZwJtSWLv0PQr2wyBiIFo9NPGZhp4T
  3P77fBaFzOTuFXVMpjYm35d4bsItn6jaue1yceSHynbPmQ3m4KiBwBj35STMGpgR95N0++
  aL6Qug/gMJdlC6oaYl17awUQAVd0vpOyJ2FS6tZy00bjNyC3///F4QKBgQDFkTDAiWKuBL
  Kk/m0Dh/4OKWjk/+QhdU1VqfsW3Yc0GWzEQhdcIg770rwW5A49fRv0PB3y7R5Im1GAbDe/
  VitYcstX9eVkd2638limakIGRHGxt4ZNxBVFwIMoyTsqD1g/y8KCq3bNffkXBwzpsXRSAw
  oxNddbdlj+ZOylt+cpTwKBgGUafRhPWlsU+jB4VW0NN/f4l9MGdeLR/Qk+PAzhWy/5dszU
  6gnO8EKflBHtTs/dBe/zZB6JN2AVoxDZcnpab8FyenmuwerNBlB1rGlogZmy0kTdoPGOdC
  svuczgCS3QdwtRmhgrHZkNcOWAoplZ2JCZ7fjJdQTO2eK/xCaYIU4hAoGAa4cDJrdfeuDZ
  rg0/nd1lO6X0Djbrnf4u8gHKw/4b+RIhbYufFSkASLTAZCEJUxQFo98YWcAObGwEZsX/bW
  bjvobz/1K43/5Ux07iSuioOKsFyjjdovOmtEj72bX0OocRe99VZTMXPO5kJNFUiNhpO72l
  zXTFWmVGOGcLmYJHEOcCgYBYp1M2CwgNYtUeJ6IGMK70cDygz7iW5vgL3LwaPlZ61CalOd
  1anFbtkRegcQ4GE7YLHYqFjzhG0trBV7CjvEZimxAkjHT51WlmwEi3M/JdSwM6YzEBFyHu
  jUe7MTYBXLh+D7eT1wYfnv/pMTUb5oItrwY+wMWx0x7AiY+C2X2yYA==
peers:
  agent:
    id: QmW1RE2V9aprXPKGXn8SqBzj34egVLaGeaVZdmZeJtthK6
```

5. Start agent service and supply the forwarding port (assume it is 8000) at agent node
```
[agent-node] $ ./p2ptunnel agent 8000
```


5. Start connector service at connector node
```
[connector-node] $ ./p2ptunnel connector
```
Connector listens at port 8012 by default. You can always change to others by `--port <your port>`
```
[connector-node] $ ./p2ptunnel connector --port <your port>
```

6. try curl at your connector node now
```
[connector-node] $ curl localhost:8012
```