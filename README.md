# consumer
consumer written in go

consumer.exe: a go program which connects to the rabbitMQ test_queue and consumes messages. When it receives a message that has a "Value" attribute which is a "1" it writes out to a text file (/tmp/stratagem_golang_output.txt) the timestamp (placed on the queue) Thu Jan 29 08:48:11 +0000 UTC 2015, the JSON object message ID 54c9f00c14339b3c3ea5052c and the string "Got a 1!". these are delimitated by the pipe symbol.

