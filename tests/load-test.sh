#!/bin/bash

vegeta attack -duration=60s -rate=100 -targets=target.list -output=attack.bin
vegeta report attack.bin
vegeta plot -title=Results attack.bin > results.html