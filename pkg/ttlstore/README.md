# Map KV Store implementation
Key-Value store with custom GC. Saves data in binary file using encoding/gob. Encoding/Decoding has been optimized, as a result file size can be up too 138x smaller. Trade of is each file can store only one TYPE of objects.
