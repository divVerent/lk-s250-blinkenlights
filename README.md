# LK-S250 Blinkenlights

A simple demo program for the proprietary SysEx messages of the
CASIO LK-S250 keyboard that control the key lights.

# The Messages

The discovered SysEx messages are as follows:

* Ping: 44 7E 7E 7F 00 03
  * This message has to be sent in regular intervals
    (below 0.5 seconds by default) in order to keep the lights on.
* LightOn: 44 7E 7E 7F 02 00 nn 01
  * This message turns on the light identified by the note number nn.
  * Up to four lights can be turned on at any given time.
* LightOff: 44 7E 7E 7F 02 00 nn 01
  * This message turns off the light identified by the note number nn.

# The Demo

The included demo program plays a simple light effects on the keyboard when
connected via USB.

# License

See the included LICENSE file.

Note: This is not an officially supported Google product.
